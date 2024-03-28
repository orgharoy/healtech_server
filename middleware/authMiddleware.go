package middleware

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Get the token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized - Missing Authorization header"})
	}

	// Check if the header has the "Bearer " prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized - Invalid Authorization header format"})
	}

	// Extract the token after removing the "Bearer " prefix
	tokenString := authHeader[7:]

	// Parse and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used for signing
		return []byte("pwoEQuF2jdk4c!nW$Nuew^rf6kjnV"), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ID := claims["ID"]
		userID := claims["userId"]
		role := claims["role"]

		c.Locals("ID", ID)
		c.Locals("userID", userID)
		c.Locals("role", role)

		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
}
