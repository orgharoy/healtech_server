package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orgharoy/healtech/database"
	helperfunctions "github.com/orgharoy/healtech/helperFunctions"
	model "github.com/orgharoy/healtech/models"
	"golang.org/x/crypto/bcrypt"
)

type EntryUser struct {
	UserName        string `gorm:"varchar(255);not null" json:"userName"`
	UserId          string `gorm:"varchar(255);not null" json:"userId"`
	Password        string `gorm:"not null" json:"password"`
	ConfirmPassword string `gorm:"not null" json:"confirmPassword"`
	Role            string `gorm:"not null" json:"role"`
}

type EntryUserForUpdate struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	Key       int64     `gorm:"int" json:"key"`
	UserName  string    `gorm:"varchar(255);not null" json:"userName"`
	UserId    string    `gorm:"varchar(255);not null" json:"userId"`
	Password  string    `gorm:"not null" json:"password"`
	Role      string    `gorm:"not null" json:"role"`
	CreatedBy string    `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
	UpdatedBy string    `gorm:"not null" json:"updatedBy"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`
	Mode      string    `gorm:"default:n" json:"mode"`
	Status    int64     `gorm:"default:1" json:"status"`
}

func CreateUser(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	db := database.DB

	var enteredUser *EntryUser

	err := c.BodyParser(&enteredUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	// -> empty fields
	if enteredUser.UserName == "" || enteredUser.UserId == "" || enteredUser.Password == "" || enteredUser.Role == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Fields cannot be empty", "data": nil})
	}

	// -> password validation

	if enteredUser.Password != enteredUser.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Passwords do not match", "data": nil})
	}

	if !helperfunctions.IsValidPassword(enteredUser.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Password must be longer than 8 characters, contain both uppercase and lowercase letter and have have special characters", "data": nil})
	}

	//find if user already exists

	var existingUser *model.User

	db.Find(&existingUser, "user_id = ?", enteredUser.UserId)

	if existingUser.ID != 0 {
		return c.Status(fiber.StatusAlreadyReported).JSON(fiber.Map{"status": "error", "message": "Account with this User ID already exists", "data": nil})
	}

	//hash password

	bytes, err := bcrypt.GenerateFromPassword([]byte(enteredUser.Password), 14)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error Hashing Password", "data": err})
	}

	//create user model from entered User details
	user := &model.User{}
	user.UserName = enteredUser.UserName
	user.UserId = enteredUser.UserId
	user.Password = string(bytes)
	user.Role = enteredUser.Role
	user.CreatedBy = c.Locals("userID").(string)
	user.Status = 2

	//create record
	err = db.Create(&user).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}

	//return created user
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "User successfully created", "data": user})
}

func GetUserForMaker(c *fiber.Ctx) error {
	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	db := database.DB

	var users []model.User
	if err := db.Order("created_at desc").Find(&users, "status = ?", 1).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to retrieve users", "data": err})
	}

	// Return the list of users with status 1
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

func GetUserForChecker(c *fiber.Ctx) error {
	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	db := database.DB

	var users []model.User
	if err := db.Order("created_at desc").Find(&users, "status = ?", 2).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to retrieve users", "data": err})
	}

	// Return the list of users with status 1
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

func ApproveUserRecord(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var updateUser EntryUserForUpdate
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	userID := c.Params("id")

	fmt.Println(userID)

	db := database.DB

	var existingUser *model.User

	db.Find(&existingUser, "id = ?", userID)

	if existingUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Account with this User ID Not Found", "data": nil})
	}

	existingUser.Status = 3
	existingUser.Mode = "u"

	if err := helperfunctions.UpdateUser(existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update user", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User status updated successfully",
		"data":    existingUser,
	})
}

func SendUserRecordToMaker(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var updateUser EntryUserForUpdate
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	userID := c.Params("id")

	fmt.Println(userID)

	db := database.DB

	var existingUser *model.User

	db.Find(&existingUser, "id = ?", userID)

	if existingUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Account with this User ID Not Found", "data": nil})
	}

	existingUser.Status = 1

	if err := helperfunctions.UpdateUser(existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update user", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User status updated successfully",
		"data":    existingUser,
	})
}

func EditUserFromActiveUserList(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var updateUser *model.User
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	userID := c.Params("id")

	db := database.DB

	var existingUser *model.User

	db.Find(&existingUser, "id = ?", userID)

	if existingUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Account with this User ID Not Found", "data": nil})
	}

	existingUser.Status = 1
	existingUser.Mode = "a"

	if err := helperfunctions.UpdateUser(existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update user", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User status updated successfully",
		"data":    existingUser,
	})
}

func GetActiveUserList(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	db := database.DB

	var activeUsers []model.User

	if err := db.Order("id desc").Find(&activeUsers, "status = ? AND mode = ?", 3, "u").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to Retrieve Active Users",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Active Users Retrieved Successfully",
		"data":    activeUsers,
	})
}

func EditUserFromMakerEnd(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var updateUser *model.User
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	userID := c.Params("id")

	db := database.DB

	var existingUser *model.User

	db.Find(&existingUser, "user_id = ?", updateUser.UserId)

	if existingUser.ID != 0 {
		return c.Status(fiber.StatusAlreadyReported).JSON(fiber.Map{"status": "error", "message": "Account with this UserID Already Exists", "data": nil})
	}

	db.Find(&existingUser, "id = ?", userID)

	if existingUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Account with this User ID Not Found", "data": nil})
	}

	existingUser.UserName = updateUser.UserName
	existingUser.UserId = updateUser.UserId
	existingUser.Role = updateUser.Role
	existingUser.UpdatedBy = c.Locals("userID").(string)
	existingUser.UpdatedAt = time.Now()
	existingUser.Status = 2

	if err := helperfunctions.UpdateUser(existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update user", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User status updated successfully",
		"data":    existingUser,
	})
}

func DeleteUser(c *fiber.Ctx) error {

	CreatedByUserRole := c.Locals("role")

	if CreatedByUserRole != "admin" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Not Authorized to Create User", "data": nil})
	}

	var updateUser *model.User
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	userID := c.Params("id")

	db := database.DB

	var existingUser *model.User

	db.Find(&existingUser, "id = ?", userID)

	if existingUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Account with this User ID Not Found", "data": nil})
	}

	existingUser.Status = 5
	existingUser.Mode = "d"

	if err := helperfunctions.UpdateUser(existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update user", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
		"data":    existingUser,
	})
}
