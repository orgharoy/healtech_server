package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/orgharoy/healtech/database"
	"github.com/orgharoy/healtech/routes"
)

func main() {
	// Start a new fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))

	err := database.Connect()

	if err != nil {
		fmt.Println(err, " failed to connect to database")
	}

	routes.Routes(app)

	app.Listen(":6969")
}
