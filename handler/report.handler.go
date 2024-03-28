package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orgharoy/healtech/database"
	helperfunctions "github.com/orgharoy/healtech/helperFunctions"
	model "github.com/orgharoy/healtech/models"
)

type EntryReport struct {
    ReportName       string `gorm:"varchar(255);not null" json:"reportName"`
    ReportGroupID    uint   `gorm:"not null" json:"reportGroupId"` // Foreign key referencing ReportGroup
    ReportGroup      model.ReportGroup `gorm:"foreignkey:ReportGroupID"` // Define relationship with ReportGroup
    ReportDescripton string `gorm:"type:text;not null" json:"reportDescription"`
    ReportPrice      string `gorm:"varchar(255);not null" json:"reportPrice"`
}

func CreateReport(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	db := database.DB

	var enteredReport *EntryReport

	err := c.BodyParser(&enteredReport)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	// -> empty fields
	if enteredReport.ReportName == "" || enteredReport.ReportDescripton == "" || enteredReport.ReportGroupID <= 0 || enteredReport.ReportPrice == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Fields cannot be empty", "data": nil})
	}

	//find if report already exists

	var existingReport *model.Report

	db.Find(&existingReport, "report_name = ?", enteredReport.ReportName)

	if existingReport.ID != 0 {
		return c.Status(fiber.StatusAlreadyReported).JSON(fiber.Map{"status": "error", "message": "Report with this name already exists", "data": nil})
	}

	//create report model from entered report details
	report := &model.Report{}

	report.ReportName = enteredReport.ReportName
	report.ReportGroupID = enteredReport.ReportGroupID
	report.ReportPrice = enteredReport.ReportPrice
	report.ReportDescripton = enteredReport.ReportDescripton
	report.CreatedBy = c.Locals("userID").(string)
	report.Status = 2

	//create record
	err = db.Create(&report).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create new report", "data": err})
	}

	//return created report
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Report successfully created", "data": report})
}

func GetReportListForMaker(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": CreatedByUserRole})
	}

	db := database.DB

	// Define a slice to store the fetched reports
	var reports []model.Report

	if err := db.Order("created_at DESC").Find(&reports, "status = ?", 1).Error; err != nil {
		// If there's an error, return an error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch reports",
			"data":    err,
		})
	}

	// Return the fetched reports in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Reports retrieved successfully",
		"data":    reports,
	})
}

func UpdateReportFromMakerEnd(c *fiber.Ctx) error {
	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var updateReport *model.Report
	if err := c.BodyParser(&updateReport); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	reportID := c.Params("id")

	db := database.DB

	var existingReport *model.Report

	db.Find(&existingReport, "report_name = ? AND id <> ?", updateReport.ReportName, reportID)

	if existingReport.ID != 0 {
		return c.Status(fiber.StatusAlreadyReported).JSON(fiber.Map{"status": "error", "message": "Report with this name already exists", "data": nil})
	}

	db.Find(&existingReport, "id = ?", reportID)

	if existingReport.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Report with this ID Not Found", "data": nil})
	}

	existingReport.ReportName = updateReport.ReportName
	existingReport.ReportGroupID = updateReport.ReportGroupID
	existingReport.ReportDescripton = updateReport.ReportDescripton
	existingReport.ReportPrice = updateReport.ReportPrice
	existingReport.UpdatedBy = c.Locals("userID").(string)
	existingReport.UpdatedAt = time.Now()
	existingReport.Status = 2

	if err := helperfunctions.UpdateReport(existingReport); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update report", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User status updated successfully",
		"data":    existingReport,
	})
}

func DeleteReport(c *fiber.Ctx) error {
	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	reportID := c.Params("id")

	db := database.DB

	var existingReport *model.Report

	db.Find(&existingReport, "id = ?", reportID)

	if existingReport.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Report with this ID Not Found", "data": nil})
	}

	existingReport.Mode = "d"
	existingReport.Status = 5

	if err := helperfunctions.UpdateReport(existingReport); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to delete report", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Report deleted successfully",
		"data":    existingReport,
	})
}

func GetReportListForChecker(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": CreatedByUserRole})
	}

	db := database.DB

	// Define a slice to store the fetched reports
	var reports []model.Report

	//if err := db.Order("created_at DESC").Find(&reports, "status = ?", 2).Error; err != nil {
	if err := db.Debug().Joins("ReportGroup").Where("reports.status = ?", 2).Order("reports.created_at DESC").Find(&reports).Error; err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch reports",
			"data":    err,
		})
	}

	// Return the fetched reports in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Reports retrieved successfully",
		"data":    reports,
	})
}

func ApproveReportRecord(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	userID := c.Params("id")

	db := database.DB

	var existingReport *model.Report

	db.Find(&existingReport, "id = ?", userID)

	if existingReport.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Report with this ID Not Found", "data": nil})
	}

	existingReport.Status = 3
	existingReport.Mode = "u"

	if err := helperfunctions.UpdateReport(existingReport); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update report", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Report status updated successfully",
		"data":    existingReport,
	})
}

func RejectReportRecord(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	userID := c.Params("id")

	db := database.DB

	var existingReport *model.Report

	db.Find(&existingReport, "id = ?", userID)

	if existingReport.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Report with this ID Not Found", "data": nil})
	}

	existingReport.Status = 1

	if err := helperfunctions.UpdateReport(existingReport); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update report", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Report successfully sent to checker",
		"data":    existingReport,
	})
}

func GetActiveReportList(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	db := database.DB

	var activeUsers []model.Report


	if err := db.Debug().Joins("ReportGroup").Where("reports.status = ? AND reports.mode = ?", 3, "u").Order("reports.created_at DESC").Find(&activeUsers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch reports",
			"data":    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Active Reports Retrieved Successfully",
		"data":    activeUsers,
	})
}

func SendReportToEditFromActive(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	userID := c.Params("id")

	db := database.DB

	var existingReport *model.Report

	db.Find(&existingReport, "id = ?", userID)

	if existingReport.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Report with this ID Not Found", "data": nil})
	}

	existingReport.Status = 1
	existingReport.Mode = "a"

	if err := helperfunctions.UpdateReport(existingReport); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update report", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Report successfully sent for ammendment",
		"data":    existingReport,
	})
}
