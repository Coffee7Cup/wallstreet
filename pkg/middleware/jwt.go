package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string
		authHeader := c.Get("Authorization")

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			// fallback to query parameter for WebSockets
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			logs.Log.Warn("Authorization failed: missing token", zap.String("path", c.Path()))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization token"})
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			logs.Log.Warn("Authorization failed: invalid or expired token", zap.String("path", c.Path()), zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", int(claims["user_id"].(float64)))
		c.Locals("username", claims["username"].(string))
		c.Locals("is_admin", claims["is_admin"].(bool))

		logs.Log.Debug("User authorized", zap.String("username", claims["username"].(string)), zap.String("path", c.Path()))

		return c.Next()
	}
}

func GenerateToken(userID int, username string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	logs.Log.Debug("Generating token", zap.String("username", username), zap.Bool("is_admin", isAdmin))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
