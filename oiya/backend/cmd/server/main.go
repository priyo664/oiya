package main

import (
	"log"
	"oiya-backend/config"
	"oiya-backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load config from .env
	config.LoadConfig()

	// Init Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Connect to DB (in config package or repository init)
	err := config.InitDB()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// Register routes
	routes.RegisterRoutes(e)

	// Start server
	port := config.AppConfig.AppPort
	log.Printf("Starting server at :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
