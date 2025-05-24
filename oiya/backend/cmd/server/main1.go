package main

import (
    "log"
    "net/http"
    "oiya-backend/config"
    "oiya-backend/routes"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    // Load konfigurasi dari .env
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Buat instance Echo
    e := echo.New()

    // Middleware logger dan recover
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Middleware CORS (boleh sesuaikan)
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
    }))

    // Inisialisasi semua route
    routes.InitRoutes(e, cfg)

    // Jalankan server
    addr := ":" + cfg.ServerPort
    log.Printf("Starting server at %s...", addr)
    if err := e.Start(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
