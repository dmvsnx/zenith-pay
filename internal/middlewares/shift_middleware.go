package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/savanyv/zenith-pay/internal/repository"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RequireActiveShift(repo repository.ShiftRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cashierID := c.Locals("userID").(string)
		cashierUUID, _ := uuid.Parse(cashierID)

		shift, err := repo.FindActiveShiftByCashier(cashierUUID.String())
		if err != nil || shift == nil {
			return helpers.Unauthorized(c, "No active shift")
		}

		c.Locals("shiftID", shift.ID.String())
		return c.Next()
	}
}
