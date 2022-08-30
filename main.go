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
	v1 := app.Group("/api/v1")
	router.NewRouter(v1)

	app.Listen(":1234")
}
