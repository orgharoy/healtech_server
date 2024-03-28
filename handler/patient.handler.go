package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/orgharoy/healtech/database"
	helperfunctions "github.com/orgharoy/healtech/helperFunctions"
	model "github.com/orgharoy/healtech/models"
)

// Report struct to represent a single report
type Report struct {
	ID                  int    `json:"id"`
	ReportName          string `json:"reportName"`
	ReportGroupId       int    `json:"reportGroupId"`
	ReportPriceCurrency string `gorm:"varchar(255);default:BDT;not null" json:"reportPriceCurrency"`
	ReportPrice         string `gorm:"varchar(255);not null" json:"reportPrice"`
	ReportDescription   string `json:"reportDescription"`
}

// Patient struct to represent the incoming JSON request body
type EnteredPatientDetails struct {
	AmountDue             string   `json:"amountDue"`
	AmountPaidReception   string   `json:"amountPaidReception"`
	AmountToPay           string   `json:"amountToPay"`
	Discount              float64  `json:"discount"`
	DiscountedAmountToPay string   `json:"discountedAmountToPay"`
	PatientAge            string   `json:"patientAge"`
	PatientContactNumber  string   `json:"patientContactNumber"`
	PatientGender         string   `json:"patientGender"`
	PatientName           string   `json:"patientName"`
	Reference             string   `json:"reference,omitempty"`
	SelectedReports       []Report `json:"selectedReports"`
	UniquePatientId       string   `json:"uniquePatientId"`
	VisitType             string   `json:"visitType"`
}

type PatientDetailsRequest struct {
	ContactWithExtension string `json:"contactWithExtension"`
	PatientID            string `json:"patientId"`
}

func NewPatientRecord(c *fiber.Ctx) error {
	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" && CreatedByUserRole != "user" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var enteredPatientDetails EnteredPatientDetails
	if err := c.BodyParser(&enteredPatientDetails); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse request body",
			"data":    nil,
		})
	}

	db := database.DB

	newPatient := &model.Patient{}
	newPatientBill := &model.PatientBill{}

	if enteredPatientDetails.VisitType == "New" {

		newPatient.UniquePatientId = helperfunctions.GenerateUniquePatientID(enteredPatientDetails.PatientName)
		newPatient.Name = enteredPatientDetails.PatientName
		newPatient.Age = enteredPatientDetails.PatientAge
		newPatient.Gender = enteredPatientDetails.PatientGender
		newPatient.ContactNumber = enteredPatientDetails.PatientContactNumber
		newPatient.CreatedBy = c.Locals("userID").(string)
		newPatient.Mode = "u"
		newPatient.Status = 3

		if err := db.Create(&newPatient).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create new patient", "data": err})
		}

	} else {

		if err := db.Find(&newPatient, "contact_number = ?", enteredPatientDetails.PatientContactNumber).Error; err != nil {
			// If there's an error, return an error response
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch reports",
				"data":    err,
			})
		}

		if newPatient.ID != 0 && newPatient.UniquePatientId != enteredPatientDetails.UniquePatientId {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Phone number already exits",
				"data":    nil,
			})
		}

		if err := db.Order("created_at DESC").Find(&newPatient, "unique_patient_id = ?", enteredPatientDetails.UniquePatientId).Error; err != nil {
			// If there's an error, return an error response
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch reports",
				"data":    err,
			})
		}

		newPatient.Name = enteredPatientDetails.PatientName
		newPatient.Age = enteredPatientDetails.PatientAge
		newPatient.Gender = enteredPatientDetails.PatientGender
		newPatient.ContactNumber = enteredPatientDetails.PatientContactNumber
		newPatient.Mode = "u"

		if err := db.Save(&newPatient).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not update existing patient", "data": err})
		}
	}

	newPatientBill.PatientID = newPatient.ID
	newPatientBill.Reference = enteredPatientDetails.Reference
	newPatientBill.UniqueBillId = helperfunctions.GenerateUniqueBillId(enteredPatientDetails.PatientName, enteredPatientDetails.VisitType)

	amountToPay, err := strconv.ParseFloat(enteredPatientDetails.AmountToPay, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid AmountToPay", "data": nil})
	}

	newPatientBill.AmountToPay = amountToPay
	newPatientBill.Discount = float32(enteredPatientDetails.Discount)

	discountedAmountToPay, err := strconv.ParseFloat(enteredPatientDetails.DiscountedAmountToPay, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid AmountToPay", "data": nil})
	}

	newPatientBill.AmountToPay = amountToPay

	newPatientBill.DiscountedAmountToPay = discountedAmountToPay

	amountPaidReception, err := strconv.ParseFloat(enteredPatientDetails.AmountPaidReception, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid AmountToPay", "data": nil})
	}

	newPatientBill.AmountPaidReception = amountPaidReception
	newPatientBill.AmountPaidFinalStage = 0

	amountDue, err := strconv.ParseFloat(enteredPatientDetails.AmountDue, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid AmountToPay", "data": nil})
	}

	newPatientBill.AmountDue = amountDue
	newPatientBill.CreatedBy = c.Locals("userID").(string)
	newPatientBill.Mode = "u"
	newPatientBill.Status = 3

	if err := db.Create(&newPatientBill).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create new patient bill", "data": err})
	}

	var billDetailsArray []Report

	for _, report := range enteredPatientDetails.SelectedReports {
		billDetail := &model.BillDetails{
			PatientBillID: newPatientBill.ID,
			ReportID:      uint(report.ID),
			CreatedBy:     c.Locals("userID").(string),
			Mode:          "u",
			Status:        3,
		}

		if err := db.Create(&billDetail).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create bill detail", "data": err})
		}

		billDetailsArray = append(billDetailsArray, report) // Append each bill detail to the array

	}

	newPatientBill.Patient = *newPatient

	response := fiber.Map{
		"status":  "success",
		"message": "New Patient Record Created",
		"data": fiber.Map{
			"patientBill": newPatientBill,
			"billDetails": billDetailsArray,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func FetchPatientDetails(c *fiber.Ctx) error {

	var requestData PatientDetailsRequest

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failed to parse request body", "data": err})
	}

	log.Println(requestData)

	var patientDetails model.Patient

	db := database.DB

	if requestData.PatientID != "" && requestData.ContactWithExtension != "" {
		if err := db.Where("unique_patient_id = ? AND contact_number  = ?", requestData.PatientID, requestData.ContactWithExtension).First(&patientDetails).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch patient details",
				"data":    nil,
			})
		}

	} else if requestData.PatientID != "" && requestData.ContactWithExtension == "" {
		if err := db.Where("unique_patient_id = ?", requestData.PatientID).First(&patientDetails).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch patient details",
				"data":    nil,
			})
		}
	} else if requestData.PatientID == "" && requestData.ContactWithExtension != "" {
		if err := db.Where("contact_number = ?", requestData.ContactWithExtension).First(&patientDetails).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch patient details",
				"data":    nil,
			})
		}

	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Patient ID or Contact Number not Provided", "data": nil})
	}

	patientDetails.ContactNumber = strings.TrimPrefix(patientDetails.ContactNumber, "+880")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Patient Details Fetched Successfully", "data": patientDetails})
}

// func PDFReportBill(c *fiber.Ctx) error {

// }
