package controllers

import (
	"encoding/json"
	"strconv"
	"time"

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

// GetProfitLoss fetches company P&L based on company ID.
func (h *MarketHandler) GetProfitLoss(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logs.Log.Warn("P&L fetch failed: invalid company id", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	plFiltered := h.store.CachePF[companyID]
	if plFiltered == nil {
		plFiltered = []models.ProfitLoss{}
	}

	return c.JSON(fiber.Map{
		"message":     "profit & loss fetched successfully",
		"profit_loss": plFiltered,
	})
}

// GetBalanceSheets fetches company Balance Sheets based on company ID.
func (h *MarketHandler) GetBalanceSheets(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logs.Log.Warn("Balance Sheets fetch failed: invalid company id", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	bsFiltered := h.store.CacheBS[companyID]
	if bsFiltered == nil {
		bsFiltered = []models.BalanceSheet{}
	}

	return c.JSON(fiber.Map{
		"message":        "balance sheets fetched successfully",
		"balance_sheets": bsFiltered,
	})
}

// GetCashFlows fetches company Cash Flows based on company ID.
func (h *MarketHandler) GetCashFlows(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logs.Log.Warn("Cash Flows fetch failed: invalid company id", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	cfFiltered := h.store.CacheCF[companyID]
	if cfFiltered == nil {
		cfFiltered = []models.CashFlow{}
	}

	return c.JSON(fiber.Map{
		"message":    "cash flows fetched successfully",
		"cash_flows": cfFiltered,
	})
}

// GetRatios fetches company ratios based on company ID.
func (h *MarketHandler) GetRatios(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logs.Log.Warn("Ratios fetch failed: invalid company id", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	ratiosFiltered := h.store.CacheRatios[companyID]
	if ratiosFiltered == nil {
		ratiosFiltered = []models.Ratio{}
	}

	return c.JSON(fiber.Map{
		"message": "ratios fetched successfully",
		"ratios":  ratiosFiltered,
	})
}

// GetStocksTillTick fetches historical stock prices for a company up to the current tick.
func (h *MarketHandler) GetStocksTillTick(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if idParam != "" {
		companyID, err := strconv.Atoi(idParam)
		if err != nil {
			logs.Log.Warn("Stock history fetch failed: invalid company id", zap.String("id", idParam))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
		}

		historyFiltered := h.store.CacheStocks[companyID]
		if historyFiltered == nil {
			historyFiltered = []models.StockPrice{}
		}

		return c.JSON(fiber.Map{
			"message": "stock history fetched successfully",
			"history": historyFiltered,
		})
	}

	allHistory := []models.StockPrice{}
	for _, companyHistory := range h.store.CacheStocks {
		allHistory = append(allHistory, companyHistory...)
	}

	return c.JSON(fiber.Map{
		"message": "stock history fetched successfully",
		"history": allHistory,
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
		allNews, err = h.store.GetNewsTillTickLimitSearchForAll(c.Context(), tickInt, limit, offset, search)
		if err != nil {
			logs.Log.Error("Failed to fetch bulk news", zap.Int("tick", tickInt), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
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

	news, err := h.store.GetNewsBySectorTillTick(c.Context(), sector, tickInt, limit, offset)
	if err != nil {
		logs.Log.Error("Failed to fetch news by sector", zap.String("sector", sector), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

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

// PeerCompanyData is the shape returned per peer in the comparison endpoint.
type PeerCompanyData struct {
	Company     models.Company       `json:"company"`
	Price       float64              `json:"price"`
	LatestPL    *models.ProfitLoss   `json:"latest_pl,omitempty"`
	LatestBS    *models.BalanceSheet `json:"latest_bs,omitempty"`
	LatestRatio *models.Ratio        `json:"latest_ratio,omitempty"`
}

// GetPeerComparison returns companies in the same sector with their calculated metrics.
func (h *MarketHandler) GetPeerComparison(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	company, err := h.store.GetCompanyByID(c.Context(), companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "company not found"})
	}

	peers, err := h.store.GetCompaniesBySector(c.Context(), company.Sector)
	if err != nil {
		logs.Log.Error("Failed to fetch peers by sector", zap.String("sector", company.Sector), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	state := h.engine.GetState()
	historyMap := h.store.CacheStocks

	var result []PeerCompanyData
	for _, peer := range peers {
		entry := PeerCompanyData{Company: peer}

		history := historyMap[peer.ID]
		var stock models.StockPrice
		var priceErr error

		if len(history) > 0 {
			stock = history[len(history)-1]
		} else {
			stock, priceErr = h.store.GetStockAtTickOfCompany(c.Context(), state.Tick, peer.ID)
		}

		price := 0.0
		var date time.Time
		if priceErr == nil {
			price = stock.ClosePrice
			date = stock.Date
		} else {
			// fallback to history if needed
			val, _ := h.store.GetStocksTillTickOfCompany(c.Context(), state.Tick, peer.ID)
			if len(val) > 0 {
				lastStock := val[len(val)-1]
				price = lastStock.ClosePrice
				date = lastStock.Date
			}
		}
		entry.Price = price

		if date.IsZero() {
			// fallback if somehow still zero
			date, _ = h.store.GetDateFromTickOfCompany(c.Context(), state.Tick, peer.ID)
		}

		if !date.IsZero() {
			pl, plErr := h.store.GetProfitLoss(c.Context(), peer.ID, date)
			if plErr == nil && len(pl) > 0 {
				entry.LatestPL = &pl[0]
			}

			bs, bsErr := h.store.GetBalanceSheets(c.Context(), peer.ID, date)
			if bsErr == nil && len(bs) > 0 {
				entry.LatestBS = &bs[0]
			}

			ratios, ratioErr := h.store.GetRatios(c.Context(), peer.ID, date)
			if ratioErr == nil && len(ratios) > 0 {
				entry.LatestRatio = &ratios[0]
			}
		}

		result = append(result, entry)
	}

	return c.JSON(fiber.Map{
		"message": "peer comparison fetched successfully",
		"peers":   result,
	})
}

// GetDrawings fetches chart drawings for a company.
func (h *MarketHandler) GetDrawings(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	drawings, err := h.store.GetDrawings(c.Context(), companyID)
	if err != nil {
		logs.Log.Error("Failed to fetch drawings", zap.Int("company_id", companyID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if drawings == nil {
		drawings = json.RawMessage("[]")
	}

	return c.JSON(fiber.Map{
		"message":  "drawings fetched successfully",
		"drawings": drawings,
	})
}

// SaveDrawings persists chart drawings for a company.
func (h *MarketHandler) SaveDrawings(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	var payload struct {
		Drawings json.RawMessage `json:"drawings"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err = h.store.SaveDrawings(c.Context(), companyID, payload.Drawings)
	if err != nil {
		logs.Log.Error("Failed to save drawings", zap.Int("company_id", companyID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message": "drawings saved successfully",
	})
}
