package helperfunctions

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/orgharoy/healtech/database"
	model "github.com/orgharoy/healtech/models"
	"github.com/unidoc/unioffice/document"
)

func IsValidPassword(password string) bool {
	// Check for at least 8 characters
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter
	uppercaseRegex := regexp.MustCompile("[A-Z]")
	if !uppercaseRegex.MatchString(password) {
		return false
	}

	// Check for at least one lowercase letter
	lowercaseRegex := regexp.MustCompile("[a-z]")
	if !lowercaseRegex.MatchString(password) {
		return false
	}

	// Check for at least one special character (you can modify the character class)
	specialCharRegex := regexp.MustCompile("[!@#$%^&*()_+]")
	return specialCharRegex.MatchString(password)
}

func UpdateUser(user *model.User) error {

	db := database.DB

	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func UpdateReport(user *model.Report) error {

	db := database.DB

	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func GenerateUniquePatientID(name string) string {
	// Get current time details
	now := time.Now()
	year := now.Year() % 100
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()

	// Get the first letter of the patient's first name
	firstLetter := strings.ToUpper(string(name[0]))

	// Generate two random letters from A-Z
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randLetter1 := charset[rand.Intn(len(charset))]
	randLetter2 := charset[rand.Intn(len(charset))]

	// Format the month and day to ensure they are two digits
	monthStr := fmt.Sprintf("%02d", month)
	dayStr := fmt.Sprintf("%02d", day)

	// Construct the unique patient ID
	uniquePatientID := fmt.Sprintf("%s%d%s%s%c%c%d%d%d", firstLetter, year, monthStr, dayStr, randLetter1, randLetter2, hour, minute, second)

	return uniquePatientID
}

func GenerateUniqueBillId(name string, visitType string) string {
	// Get current time details
	now := time.Now()
	year := now.Year() % 100
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()

	// Get the first letter of the patient's first name
	firstLetter := strings.ToUpper(string(name[0]))
	visitTypeLetter := strings.ToUpper(string(visitType[0]))

	// Generate two random letters from A-Z
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randLetter1 := charset[rand.Intn(len(charset))]
	randLetter2 := charset[rand.Intn(len(charset))]
	randLetter3 := charset[rand.Intn(len(charset))]

	// Format the month and day to ensure they are two digits
	monthStr := fmt.Sprintf("%02d", month)
	dayStr := fmt.Sprintf("%02d", day)

	// Construct the unique patient ID
	uniqueBillId := fmt.Sprintf("%s%c%d%d%d%s%d%s%s%c%c", visitTypeLetter, randLetter3, hour, minute, second, firstLetter, year, monthStr, dayStr, randLetter1, randLetter2)
	//uniquePatientID := fmt.Sprintf("%s%s%d%s%s%c%c%d%d%d", visitTypeLetter, firstLetter, year, monthStr, dayStr, randLetter1, randLetter2, hour, minute, second)

	return uniqueBillId
}

func ReplacePlaceholder(doc *document.Document, placeholder, value string) {
    for _, p := range doc.Paragraphs() {
        for _, r := range p.Runs() {
            txt := r.Text()
            if strings.Contains(txt, placeholder) {
                txt = strings.Replace(txt, placeholder, value, -1)
                r.Clear()
                r.AddText(txt)
            }
        }
    }
}

func ConvertToPDF(inputFile, outputFile string) error {
    // unoconv command to convert Word to PDF
    cmd := exec.Command("unoconv", "-f", "pdf", inputFile)

    // Set output file for PDF
    pdfFile, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer pdfFile.Close()
    cmd.Stdout = pdfFile

    // Execute the command
    if err := cmd.Run(); err != nil {
        return err
    }

    return nil
}