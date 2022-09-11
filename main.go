package main

import (
	"log"
	"pawang-backend/config"
	"pawang-backend/router"
	"pawang-backend/seeder"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Database
	db, err := config.Database()
	seeder.SeederCategory(db)

	if err != nil {
		log.Fatal(err.Error())
	}

	// API VERSIONING
	api := app.Group("/api")
	router.NewRouter(api)

	app.Listen(":1234")
}
