package middleware

import (
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

var ActiveRequests int32

// RequestTracker increments a counter when a request starts and decrements when it ends.
func RequestTracker() fiber.Handler {
	return func(c *fiber.Ctx) error {
		atomic.AddInt32(&ActiveRequests, 1)
		defer atomic.AddInt32(&ActiveRequests, -1)
		return c.Next()
	}
}

func GetActiveRequests() int32 {
	return atomic.LoadInt32(&ActiveRequests)
}
