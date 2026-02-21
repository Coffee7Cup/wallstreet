package routers

import (
	"github.com/coffee7cup/wallstreet/pkg/controllers"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func MarketRouter(router fiber.Router, handler *controllers.MarketHandler) {
	logs.Log.Info("Initializing Market Router")

	router.Get("/companies", handler.GetAllCompanies)
	router.Get("/sectors", handler.GetAllSectors)
	router.Get("/companies/:id", handler.GetCompany)
	router.Get("/fundamentals/:id", handler.GetFundamentals)
	router.Get("/ratios/:id", handler.GetRatios)
	router.Get("/stocks-history/:id", handler.GetStocksTillTick)
	router.Get("/stocks-history", handler.GetStocksTillTick)

	router.Get("/news", middleware.JWTMiddleware(), handler.GetNewsTillTickLimitSearch)
	router.Get("/news/sector", middleware.JWTMiddleware(), handler.GetNewsBySectorTillTick)
}
