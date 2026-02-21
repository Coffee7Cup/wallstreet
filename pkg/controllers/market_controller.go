package controllers

import (
	"strconv"

	"sort"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/coffee7cup/wallstreet/pkg/simulation"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type MarketHandler struct {
	store  *db.Store
	engine *simulation.Engine
}

func NewMarketHandler(store *db.Store, engine *simulation.Engine) *MarketHandler {
	return &MarketHandler{store: store, engine: engine}
}

// GetAllCompanies fetches all companies.
func (h *MarketHandler) GetAllCompanies(c *fiber.Ctx) error {
	companies, err := h.store.GetAllCompanies(c.Context())
	if err != nil {
		logs.Log.Error("Failed to fetch companies", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message":   "companies fetched successfully",
		"companies": companies,
	})
}

// GetCompany fetches a single company by ID.
func (h *MarketHandler) GetCompany(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logs.Log.Warn("Company fetch failed: invalid company id", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	company, err := h.store.GetCompanyByID(c.Context(), companyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "company not found"})
		}
		logs.Log.Error("Company fetch failed: internal error", zap.Int("company_id", companyID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message": "company fetched successfully",
		"company": company,
	})
}

// GetFundamentals fetches company fundamentals based on company ID.
func (h *MarketHandler) GetFundamentals(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logs.Log.Warn("Fundamentals fetch failed: invalid company id", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	fundamentals, err := h.store.GetFundamentals(c.Context(), companyID)
	if err != nil {
		logs.Log.Error("Fundamentals fetch failed: internal error", zap.Int("company_id", companyID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message":      "fundamentals fetched successfully",
		"fundamentals": fundamentals,
	})
}

// GetRatios fetches company ratios based on company ID.
func (h *MarketHandler) GetRatios(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logs.Log.Warn("Ratios fetch failed: invalid company id", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	ratios, err := h.store.GetRatios(c.Context(), companyID)
	if err != nil {
		logs.Log.Error("Ratios fetch failed: internal error", zap.Int("company_id", companyID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message": "ratios fetched successfully",
		"ratios":  ratios,
	})
}

// GetStocksTillTick fetches historical stock prices for a company up to the current tick.
func (h *MarketHandler) GetStocksTillTick(c *fiber.Ctx) error {
	idParam := c.Params("id")
	state := h.engine.GetState()

	if idParam != "" {
		companyID, err := strconv.Atoi(idParam)
		if err != nil {
			logs.Log.Warn("Stock history fetch failed: invalid company id", zap.String("id", idParam))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
		}

		history, err := h.store.GetStocksTillTickOfCompany(c.Context(), state.Tick, companyID)
		if err != nil {
			logs.Log.Error("Stock history fetch failed: internal error", zap.Int("company_id", companyID), zap.Int("tick", state.Tick), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		return c.JSON(fiber.Map{
			"message": "stock history fetched successfully",
			"history": history,
		})
	}

	history, err := h.store.GetStocksTillTick(c.Context(), state.Tick)
	if err != nil {
		logs.Log.Error("Stock history fetch failed: internal error", zap.Int("tick", state.Tick), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message": "stock history fetched successfully",
		"history": history,
	})
}

func (h *MarketHandler) GetNewsTillTickLimitSearch(c *fiber.Ctx) error {
	state := h.engine.GetState()
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	search := c.Query("search")
	tick := c.Query("tick")
	companyIDStr := c.Query("company_id")

	var tickInt int
	var err error
	if tick == "" {
		tickInt = state.Tick
	} else {
		tickInt, err = strconv.Atoi(tick)
		if err != nil || tickInt > state.Tick {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tick"})
		}
	}

	var allNews []models.News
	if companyIDStr != "" {
		companyID, err := strconv.Atoi(companyIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company_id"})
		}
		date, err := h.store.GetDateFromTickOfCompany(c.Context(), tickInt, companyID)
		if err == nil {
			allNews, err = h.store.GetNewsSearchLimitOffsetByCompany(c.Context(), companyID, date, limit, offset, search)
			if err != nil {
				logs.Log.Error("Failed to fetch company-specific news", zap.Int("company_id", companyID), zap.Error(err))
			}
		}
	} else {
		companies, err := h.store.GetAllCompanies(c.Context())
		if err != nil {
			logs.Log.Error("Failed to fetch companies for news", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		for _, company := range companies {
			date, err := h.store.GetDateFromTickOfCompany(c.Context(), tickInt, company.ID)
			if err != nil {
				continue
			}

			newsList, err := h.store.GetNewsSearchLimitOffsetByCompany(c.Context(), company.ID, date, limit+offset, 0, search)
			if err != nil {
				continue
			}
			allNews = append(allNews, newsList...)
		}

		// Sort allNews by date descending
		sort.Slice(allNews, func(i, j int) bool {
			return allNews[i].ReleaseDate.After(allNews[j].ReleaseDate)
		})

		// Apply pagination after sorting
		start := offset
		if start > len(allNews) {
			start = len(allNews)
		}
		end := start + limit
		if end > len(allNews) {
			end = len(allNews)
		}
		allNews = allNews[start:end]
	}

	return c.JSON(fiber.Map{
		"message": "news fetched successfully",
		"news":    allNews,
		"tick":    tickInt,
	})
}

func (h *MarketHandler) GetNewsBySectorTillTick(c *fiber.Ctx) error {
	state := h.engine.GetState()
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	tick := c.Query("tick")
	sector := c.Query("sector")

	tickInt, _ := strconv.Atoi(tick)
	if state.Tick < tickInt {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "you have give invalid tick i.e tick > current tick"})
	}

	companies, err := h.store.GetAllCompanies(c.Context())
	if err != nil {
		logs.Log.Error("Failed to fetch companies for news", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	var sectorCompanies []models.Company
	for _, company := range companies {
		if company.Sector == sector {
			sectorCompanies = append(sectorCompanies, company)
		}
	}

	var news []models.News

	for _, company := range sectorCompanies {
		date, err := h.store.GetDateFromTickOfCompany(c.Context(), tickInt, company.ID)
		if err != nil {
			logs.Log.Error("Failed to fetch dates for the companies", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		newsList, err := h.store.GetNewsSearchLimitOffsetByCompany(c.Context(), company.ID, date, limit+offset, 0, "")
		if err != nil {
			logs.Log.Error("Failed to fetch news for company", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		news = append(news, newsList...)
	}

	// Sort allNews by date descending
	sort.Slice(news, func(i, j int) bool {
		return news[i].ReleaseDate.After(news[j].ReleaseDate)
	})

	return c.JSON(fiber.Map{
		"message": "news fetched successfully",
		"news":    news,
		"tick":    tickInt,
	})
}

func (h *MarketHandler) GetAllSectors(c *fiber.Ctx) error {
	sectors, err := h.store.GetAllSectors(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.JSON(fiber.Map{
		"message": "sectors fetched successfully",
		"sectors": sectors,
	})
}
