package main

import (
	"pawang-backend/config"
	"pawang-backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	// Connect Database
	config.ConnectDatabase()

	// Setup Router
	routes.SetupRouter(e)

	e.Logger.Fatal(e.Start(":1234"))
}
