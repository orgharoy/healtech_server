package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orgharoy/healtech/database"
	model "github.com/orgharoy/healtech/models"
)

func GetActiveReportGroupList(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized", "data": nil})
	}

	db := database.DB

	var activeReportGroup []model.ReportGroup

	if err := db.Order("created_at desc").Find(&activeReportGroup, "status = ? AND mode = ?", 3, "u").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to Retrieve Records",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Activ Report Groups Retrieved Successfully",
		"data":    activeReportGroup,
	})
}