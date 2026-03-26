package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RoleMiddleware(allowedRoles ...model.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleStr, ok := c.Locals("role").(string)
		if !ok {
			return helpers.Forbidden(c, "access denied")
		}
		userRoleStr = strings.ToLower(userRoleStr)

		for _, role := range allowedRoles {
			if userRoleStr == string(role) {
				return c.Next()
			}
		}

		return helpers.Forbidden(c, "forbidden: infufficient permissions")
	}
}
