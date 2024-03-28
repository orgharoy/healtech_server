package handler

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/orgharoy/healtech/database"
	model "github.com/orgharoy/healtech/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	UserId   string `gorm:"not null" json:"userId"`
	Password string `gorm:"not null" json:"password"`
}

func Login(c *fiber.Ctx) error {

	db := database.DB

	userLogin := new(LoginInput)

	err := c.BodyParser(userLogin)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid Email/Username or Pasword", "data": err})
	}

	// -> empty fields
	if userLogin.UserId == "" || userLogin.Password == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Fields Cannot be Empty", "data": nil})
	}

	var existingUser model.User

	db.Where("user_id = ? AND status = ? AND mode = ?", userLogin.UserId, 3, "u").Find(&existingUser)

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userLogin.Password))

	if err != nil {
		fmt.Println(userLogin)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid Username or Password!", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()

	// Set claims (payload data) in the token.
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = existingUser.ID
	claims["userId"] = existingUser.UserId
	claims["role"] = existingUser.Role
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(time.Minute * 15).Unix()

	// Sign the token with the secret key.
	tokenString, err := token.SignedString([]byte("pwoEQuF2jdk4c!nW$Nuew^rf6kjnV"))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "generating JWT Token failed"})
	}
	// send the user's model as well as the token.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Login Successful", "token": tokenString, "user": existingUser})

}
