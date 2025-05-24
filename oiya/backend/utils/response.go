package utils

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

// ResponseSuccess mengirim response sukses dengan data
func ResponseSuccess(c echo.Context, message string, data interface{}) error {
    return c.JSON(http.StatusOK, map[string]interface{}{
        "status":  "success",
        "message": message,
        "data":    data,
    })
}

// ResponseCreated mengirim response saat data berhasil dibuat
func ResponseCreated(c echo.Context, message string, data interface{}) error {
    return c.JSON(http.StatusCreated, map[string]interface{}{
        "status":  "success",
        "message": message,
        "data":    data,
    })
}

// ResponseError mengirim response error umum
func ResponseError(c echo.Context, statusCode int, message string) error {
    return c.JSON(statusCode, map[string]interface{}{
        "status":  "error",
        "message": message,
    })
}

// ResponseValidationError mengirim response jika validasi gagal
func ResponseValidationError(c echo.Context, errors interface{}) error {
    return c.JSON(http.StatusBadRequest, map[string]interface{}{
        "status":  "error",
        "message": "validation failed",
        "errors":  errors,
    })
}

// UnauthorizedResponse untuk token invalid
func UnauthorizedResponse(c echo.Context) error {
    return c.JSON(http.StatusUnauthorized, map[string]interface{}{
        "status":  "error",
        "message": "unauthorized",
    })
}
