package controllers

import (
    "database/sql"
    "net/http"
    "github.com/labstack/echo/v4"
)

func Register(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"message": "Register OK"})
    }
}

func Login(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"message": "Login OK"})
    }
}
