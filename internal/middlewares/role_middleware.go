package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RoleMiddleware(allowedRoles ...model.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleValue := c.Locals("role")
		if roleValue == nil {
			return helpers.Forbidden(c, "you don't have permission to access this resource")
		}

		userRoleStr, ok := roleValue.(string)
		if !ok {
			return helpers.Forbidden(c, "you don't have permission to access this resource")
		}

		userRole := model.Role(userRoleStr)
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return helpers.Forbidden(c, "you don't have permission to access this resource")
	}
}
