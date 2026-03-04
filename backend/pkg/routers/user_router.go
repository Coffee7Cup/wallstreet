package routers

import (
	"github.com/coffee7cup/wallstreet/pkg/controllers"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router, handler *controllers.AuthHandler) {
	logs.Log.Info("Initializing User Router")
	router.Post("/login", handler.UserLoginHandler)
	router.Get("/profile", middleware.JWTMiddleware(), handler.GetUserProfile)
}
