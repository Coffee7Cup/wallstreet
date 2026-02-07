package controllers

import (
	"strconv"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{store: store}
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
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
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
