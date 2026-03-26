package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func JWTMiddleware(jwtService helpers.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if authHeader == "" {
			return helpers.Unauthorized(c, "Missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return helpers.Unauthorized(c, "Invalid Authorization format")
		}
		tokenString := parts[1]

		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil || claims == nil {
			return helpers.Unauthorized(c, "Invalid or expired JWT")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("claims", claims)

		return c.Next()
	}
}
