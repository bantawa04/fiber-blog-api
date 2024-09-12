package utils

import (
	"github.com/gofiber/fiber/v2"
	"math"
)

type ResponseUtil struct{}

type PaginationMeta struct {
	Page       int   `json:"page"`
	TotalPages int   `json:"total_pages"`
	PerPage    int   `json:"per_page"`
	TotalItems int64 `json:"total_items"`
}

func (ru *ResponseUtil) SendResponse(c *fiber.Ctx, data interface{}, message string) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func (ru *ResponseUtil) SendError(c *fiber.Ctx, message string, description string, code int) error {
	return c.Status(code).JSON(fiber.Map{
		"success":     false,
		"message":     message,
		"description": description,
	})
}

func (ru *ResponseUtil) SendSuccess(c *fiber.Ctx, message string) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}

func (ru *ResponseUtil) SendPagination(c *fiber.Ctx, data interface{}, total int64, page, limit int, message string) error {
	meta := PaginationMeta{
		Page:       page,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		PerPage:    limit,
		TotalItems: total,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
