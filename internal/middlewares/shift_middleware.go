package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/savanyv/zenith-pay/internal/repository"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RequireActiveShift(repo repository.ShiftRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cashierID, ok := c.Locals("userID").(string)
		if !ok {
			return helpers.Unauthorized(c, "Unauthorized")
		}

		cashierUUID, err := uuid.Parse(cashierID)
		if err != nil {
			return helpers.Unauthorized(c, "Unauthorized")
		}

		shift, err := repo.FindActiveShiftByCashier(cashierUUID.String())
		if err != nil {
			return helpers.InternalServerError(c, "failed to check shift")
		}

		if shift == nil {
			return helpers.Forbidden(c, "No active shift")
		}

		c.Locals("shiftID", shift.ID.String())

		return c.Next()
	}
}
