package main

import (
	"pawang-backend/config"
	"pawang-backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func main() {
	goth.UseProviders(google.New(
		config.GetEnv("GOOGLE_CLIENT_ID"),
		config.GetEnv("GOOGLE_CLIENT_SECRET"),
		config.GetEnv("GOOGLE_CALLBACK_REDIRECT"),
	))

	e := echo.New()

	e.Use(middleware.Logger())

	// Connect Database
	config.ConnectDatabase()

	// Setup Router
	routes.SetupRouter(e)

	e.Logger.Fatal(e.Start(":1234"))
}
