package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func MethodValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		allowedMethods := map[string]bool{
			http.MethodGet:     true,
			http.MethodPost:    true,
			http.MethodPut:     true,
			http.MethodPatch:   true,
			http.MethodDelete:  true,
			http.MethodOptions: true,
		}

		if !allowedMethods[c.Method()] {
			return helpers.MethodNotAllowed(c, "Method not allowed")
		}

		return c.Next()
	}
}
