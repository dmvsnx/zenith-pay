package helpers

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Code    int         `json:"code,omitempty"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *fiber.Ctx, statusCode int, status string, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Code:    statusCode,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusOK, "success", message, data)
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusCreated, "success", message, data)
}

func BadRequest(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusBadRequest, "error", message, nil)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusUnauthorized, "error", message, nil)
}

func NotFound(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusNotFound, "error", message, nil)
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusInternalServerError, "error", message, nil)
}

func TooManyRequests(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusTooManyRequests, "error", message, nil)
}

func MethodNotAllowed(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusMethodNotAllowed, "error", message, nil)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusForbidden, "error", message, nil)
}
