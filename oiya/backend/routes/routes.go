package routes

import (
    "database/sql"
    "github.com/labstack/echo/v4"
    "oiya/controllers"
)

func SetupRoutes(e *echo.Echo, db *sql.DB) {
    e.POST("/api/register", controllers.Register(db))
    e.POST("/api/login", controllers.Login(db))
}
