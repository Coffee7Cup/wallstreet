package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/coffee7cup/wallstreet/pkg/controllers"
	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/coffee7cup/wallstreet/pkg/routers"
	"github.com/coffee7cup/wallstreet/pkg/simulation"
	"go.uber.org/zap"

	"github.com/coffee7cup/wallstreet/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	logs.Init()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			logs.Log.Error("HTTP Error", zap.Error(err), zap.String("path", c.Path()))
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(middleware.RequestTracker())

	logs.Log.Info("Starting server...")

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	err := godotenv.Load()
	if err != nil {
		logs.Log.Fatal("Failed to load environment variables", zap.Error(err))
	}

	DATABASE_URL := os.Getenv("DB_URL")

	store, err := db.NewStore(DATABASE_URL)
	if err != nil {
		logs.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer store.Close()

	engine := simulation.NewEngine(store)
	hub := utils.NewHub(engine.Subscribe(), store, engine)

	// Start engine immediately for development convenience
	engine.Start(0)

	// Root context for background tasks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go engine.Run(ctx)
	go hub.Run(ctx)

	userGroup := app.Group("/api/v1/users")
	authHandler := controllers.NewAuthHandler(store)
	routers.UserRouter(userGroup, authHandler)

	adminGroup := app.Group("/api/v1/admin")
	routers.AdminRouter(adminGroup, controllers.NewAdminHandler(store, engine, hub))

	tradeGroup := app.Group("/api/v1/trade")
	tradeHandler := controllers.NewTradeHandler(store, engine, hub)
	routers.TradeRouter(tradeGroup, tradeHandler)

	marketGroup := app.Group("/api/v1/market")
	marketHandler := controllers.NewMarketHandler(store, engine)
	routers.MarketRouter(marketGroup, marketHandler)

	logs.Log.Info("Server listening on :3000")

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			logs.Log.Error("Server error", zap.Error(err))
		}
	}()

	// Wait for termination signal
	<-stop
	logs.Log.Info("Shutting down server...")
	cancel() // Stop background tasks

	// 1. Shutdown Fiber
	if err := app.Shutdown(); err != nil {
		logs.Log.Error("Fiber shutdown error", zap.Error(err))
	}

	// 2. Stop Simulation Engine
	engine.Stop()
	logs.Log.Info("Simulation engine stopped")

	// 3. Close Database Store
	store.Close()
	logs.Log.Info("Database connection closed")

	logs.Log.Info("Goodbye!")
}
