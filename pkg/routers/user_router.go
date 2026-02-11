package routers

import (
	"github.com/coffee7cup/wallstreet/pkg/controllers"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router, handler *controllers.UserHandler) {
	logs.Log.Info("Initializing User Router")
	router.Post("/login", handler.UserLoginHandler)
	router.Get("/profile", middleware.JWTMiddleware(), handler.GetUserProfile)
	router.Get("/portfolio", middleware.JWTMiddleware(), handler.GetUserPortfolio)

	//the below routes need a refactoring
	router.Get("/trades", middleware.JWTMiddleware(), handler.GetUserTrades)
	router.Get("/companies", handler.GetAllCompanies)
	router.Get("/companies/:id", handler.GetCompany)
	router.Get("/fundamentals/:id", handler.GetFundamentals)
	router.Get("/ratios/:id", handler.GetRatios)
	router.Get("/stocks-history/:id", handler.GetStocksTillTick)
	router.Get("/stocks-history", handler.GetStocksTillTick)
	
	router.Get("/news", middleware.JWTMiddleware(), handler.GetNewsTillTick)
	router.Get("/news/:company_id", middleware.JWTMiddleware(), handler.GetNewsTillTickOfCompany)

}
