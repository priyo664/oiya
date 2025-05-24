package controllers

import (
    "database/sql"
    "net/http"

    "github.com/labstack/echo/v4"
    "oiya/middleware"
)

type Trip struct {
    Origin      string `json:"origin"`
    Destination string `json:"destination"`
}

func CreateTrip(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        userID := middleware.GetUserID(c)

        var trip Trip
        if err := c.Bind(&trip); err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trip data"})
        }

        _, err := db.Exec("INSERT INTO trips (user_id, origin, destination) VALUES (?, ?, ?)",
            userID, trip.Origin, trip.Destination)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create trip"})
        }

        return c.JSON(http.StatusOK, map[string]string{"message": "Trip created"})
    }
}
