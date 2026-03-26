package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RateLimiter(max int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max: max,
		Expiration: duration,

		KeyGenerator: func(c *fiber.Ctx) string {
			if userID, ok := c.Locals("userID").(string); ok {
				return userID
			}

			return c.IP() + c.Get("User-Agent")
		},

		LimitReached: func(c *fiber.Ctx) error {
			return helpers.TooManyRequests(c, "Too many requests, please try again later.")
		},
	})
}
