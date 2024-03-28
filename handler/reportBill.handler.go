package handler

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orgharoy/healtech/database"
	helperfunctions "github.com/orgharoy/healtech/helperFunctions"
	model "github.com/orgharoy/healtech/models"
	"github.com/unidoc/unioffice/document"
	"gorm.io/gorm"
)

type EntryBillStruct struct {

	AmountDue string `gorm:"not null" json:"amountDue"`
	AmountPaidFinalStage float64 `gorm:"not null" json:"AmountPaidFinalStage"`
	AmountPaidReception float64 `gorm:"not null" json:"amountPaidReception"`
	AmountToPay float64 `gorm:"not null" json:"amountToPay"`
	BillNumber string `gorm:"varchar(255);not null" json:"billNumber"`
	Discount float64 `gorm:"not null" json:"discount"`
	DiscountedAmountToPay float64 `gorm:"not null" json:"discountedAmountToPay"`
}

func FetchBillDetails(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" && CreatedByUserRole != "user" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	billNumber := c.Params("billId")

	db := database.DB

	var existingBill *model.PatientBill

	if err := db.Debug().Joins("Patient").Where("patient_bills.unique_bill_id = ?", strings.ToUpper(billNumber)).Find(&existingBill).Error; err != nil {

		if db.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Bill not found", "data": nil})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch reports",
			"data":    err,
		})
	}

	existingBill.Patient.ContactNumber = strings.TrimPrefix(existingBill.Patient.ContactNumber, "+880")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Bill Details Fetched Successfully",
		"data":    existingBill,
	})
}

func UpdateReportBill(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var updateBill EntryBillStruct
	if err := c.BodyParser(&updateBill); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	db := database.DB

	var existingBill *model.PatientBill

	db.Find(&existingBill, "unique_bill_id = ?", updateBill.BillNumber)

	if existingBill.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Bill with this ID Not Found", "data": nil})
	}

	existingBill.AmountPaidFinalStage = updateBill.AmountPaidFinalStage
	amountDueConverted, err := strconv.ParseFloat(updateBill.AmountDue, 64)
	if err != nil {
		// Handle the error if conversion fails
		return err
	}

	// Assign the converted value to existingBill.AmountDue
	existingBill.AmountDue = amountDueConverted
	existingBill.UpdatedBy = c.Locals("userID").(string)
	existingBill.UpdatedAt = time.Now()

	if err := database.DB.Save(&existingBill).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update bill", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User status updated successfully",
		"data":    existingBill,
	})
}

func GetPDFReport(c *fiber.Ctx) error {

	doc, err := document.Open("report/PatientReportTemplate.docx")

	if err != nil {
		return err
	}
	
	helperfunctions.ReplacePlaceholder(doc, "{{PATIENT_NAME}}", "John Doe")

	if err = doc.SaveToFile("output.docx"); err != nil{
		panic(err)
	}

	defer doc.Close()

	if err := helperfunctions.ConvertToPDF("output.docx", "output.pdf"); err != nil {
        return err
    }

	pdfData, err := os.ReadFile("output.pdf")
	if err != nil {
		return err
	}

	os.Remove("output.docx")
	os.Remove("output.pdf")

	c.Set(fiber.HeaderContentType, "application/pdf")
	c.Set(fiber.HeaderContentDisposition, "attachment; filename=report.pdf")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "PDF Created Successfully",
		"data":    pdfData,
	})
}

