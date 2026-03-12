package middlewares

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func JWTMiddleware(jwtService helpers.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helpers.Unauthorized(c, "Missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helpers.Unauthorized(c, "Missing or malformed JWT")
		}
		tokenString := parts[1]

		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			return helpers.Unauthorized(c, "Invalid or expired JWT")
		}

		if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
			return helpers.Unauthorized(c, "Token expired")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("claims", claims)

		return c.Next()
	}
}
