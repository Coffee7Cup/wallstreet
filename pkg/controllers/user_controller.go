package controllers

import (
	"strconv"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/coffee7cup/wallstreet/pkg/simulation"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type UserHandler struct {
	store  *db.Store
	engine *simulation.Engine
}

func NewUserHandler(store *db.Store, engine *simulation.Engine) *UserHandler {
	return &UserHandler{store: store, engine: engine}
}

// UserLoginHandler handles user authentication.
func (h *UserHandler) UserLoginHandler(c *fiber.Ctx) error {

	type Body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var body Body

	if err := c.BodyParser(&body); err != nil {
		logs.Log.Warn("User login failed: invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if body.Username == "" && body.Email == "" {
		logs.Log.Warn("User login failed: missing username or email")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username or email is required",
		})
	}

	if body.Username != "" && body.Email != "" {
		logs.Log.Warn("User login failed: both username and email provided")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "provide either username or email, not both",
		})
	}

	logs.Log.Info("User login attempt", zap.String("username", body.Username), zap.String("email", body.Email))

	// fetch user
	user, err := h.store.GetUser(c.Context(), body.Username, body.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			logs.Log.Warn("User login failed: invalid credentials", zap.String("username", body.Username))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}

		logs.Log.Error("User login failed: internal error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, false)
	if err != nil {
		logs.Log.Error("User login failed: token generation error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	logs.Log.Info("User login successful", zap.String("username", user.Username))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login successful",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"token": token,
	})
}

// GetUserProfile fetches the profile of the authenticated user.
func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	username, ok := c.Locals("username").(string)
	if !ok {
		logs.Log.Warn("Profile fetch failed: unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user, err := h.store.GetUser(c.Context(), username, "")
	if err != nil {
		if err == pgx.ErrNoRows {
			logs.Log.Warn("Profile fetch failed: user not found", zap.String("username", username))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}

		logs.Log.Error("Profile fetch failed: internal error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user profile fetched successfully",
		"user": fiber.Map{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"cash_balance": user.CashBalance,
		},
	})
}

// GetUserPortfolio fetches the stock portfolio of the authenticated user.
func (h *UserHandler) GetUserPortfolio(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(int)
	if !ok {
		logs.Log.Warn("Portfolio fetch failed: unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	portfolio, err := h.store.GetUserPortfolio(c.Context(), id)
	if err != nil {
		logs.Log.Error("Portfolio fetch failed: internal error", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
			"err":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "user portfolio fetched successfully",
		"portfolio": portfolio,
	})
}

// GetUserTrades fetches the trade history of the authenticated user.
func (h *UserHandler) GetUserTrades(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(int)
	if !ok {
		logs.Log.Warn("Trades fetch failed: unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	trades, err := h.store.GetUserTrades(c.Context(), id)
	if err != nil {
		logs.Log.Error("Trades fetch failed: internal error", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user trades fetched successfully",
		"trades":  trades,
	})
}

// GetFundamentals fetches company fundamentals based on company ID.
func (h *UserHandler) GetFundamentals(c *fiber.Ctx) error {
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
func (h *UserHandler) GetRatios(c *fiber.Ctx) error {
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

// GetAllCompanies fetches all companies.
func (h *UserHandler) GetAllCompanies(c *fiber.Ctx) error {
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
func (h *UserHandler) GetCompany(c *fiber.Ctx) error {
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

// GetStocksTillTick fetches historical stock prices for a company up to the current tick.
func (h *UserHandler) GetStocksTillTick(c *fiber.Ctx) error {
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

// GetNewsTillTick fetches news for all companies up to the current tick.
func (h *UserHandler) GetNewsTillTickOfCompany(c *fiber.Ctx) error {
	state := h.engine.GetState()
	companyID, err := strconv.Atoi(c.Params("company_id"))
	if err != nil {
		logs.Log.Warn("News fetch failed: invalid company id", zap.String("id", c.Params("company_id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid company id"})
	}

	date, err := h.store.GetDateFromTickOfCompany(c.Context(), state.Tick, companyID)
	if err != nil {
		logs.Log.Error("News fetch failed: internal error", zap.Int("tick", state.Tick), zap.Int("company_id", companyID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	news, err := h.store.GetNewsTillDateByCompanyId(c.Context(), date, companyID)
	if err != nil {
		logs.Log.Error("News fetch failed: internal error", zap.Time("date", date), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"message": "news fetched successfully",
		"news":    news,
	})
}


func (h *UserHandler) GetNewsTillTick(c *fiber.Ctx) error {
	state := h.engine.GetState()

	companies, err := h.store.GetAllCompanies(c.Context())
	if err != nil {
		logs.Log.Error("Failed to fetch companies", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	var news []models.News
	for _, company := range companies {

		date, err := h.store.GetDateFromTickOfCompany(c.Context(), state.Tick, company.ID)
		if err != nil {
			logs.Log.Error("News fetch failed: internal error", zap.Int("tick", state.Tick), zap.Int("company_id", company.ID), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		news_list, err := h.store.GetNewsTillDateByCompanyId(c.Context(), date, company.ID)
		if err != nil {
			logs.Log.Error("News fetch failed: internal error", zap.Time("date", date), zap.Int("company_id", company.ID), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		news = append(news, news_list...)
	}

	return c.JSON(fiber.Map{
		"message": "news fetched successfully",
		"news":    news,
	})
}