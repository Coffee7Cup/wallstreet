package middleware

import (
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func VerifyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		username := c.Locals("username")
		if !ok || !isAdmin {
			logs.Log.Warn("Access denied: admin privileges required", zap.Any("username", username), zap.String("path", c.Path()))
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "access denied: admin privileges required",
			})
		}
		return c.Next()
	}
}
