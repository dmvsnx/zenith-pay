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

type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedAPIResponse struct {
	Code       int             `json:"code,omitempty"`
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Pagination *PaginationMeta `json:"pagination"`
	Data       interface{}     `json:"data,omitempty"`
}

func PaginatedSuccess(c *fiber.Ctx, message string, items interface{}, total int64, page, limit int) error {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	return c.Status(fiber.StatusOK).JSON(PaginatedAPIResponse{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: message,
		Pagination: &PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
		Data: items,
	})
}
