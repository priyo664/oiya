package controllers

import (
    "net/http"
    "oiya-backend/models"
    "oiya-backend/utils"
    "time"

    "github.com/labstack/echo/v4"
)

// RegisterUser untuk register penumpang atau driver (sederhana)
func RegisterUser(c echo.Context) error {
    req := new(models.UserRegisterRequest)
    if err := c.Bind(req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request"))
    }

    // Validasi data dasar (contoh)
    if req.Phone == "" || req.Password == "" {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Phone and Password required"))
    }

    // Cek user sudah ada di DB (pseudo code)
    exists, err := models.CheckUserExists(req.Phone)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("DB error"))
    }
    if exists {
        return c.JSON(http.StatusConflict, utils.NewErrorResponse("User already registered"))
    }

    // Hash password
    hashedPass, err := utils.HashPassword(req.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to hash password"))
    }

    // Simpan user ke DB (pseudo code)
    user := models.User{
        Phone:    req.Phone,
        Password: hashedPass,
        Role:     req.Role, // "passenger" atau "driver"
        CreatedAt: time.Now(),
    }
    if err := models.CreateUser(&user); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create user"))
    }

    return c.JSON(http.StatusCreated, utils.NewSuccessResponse("User registered successfully"))
}

// LoginUser untuk login dengan phone + OTP (sederhana tanpa OTP service)
func LoginUser(c echo.Context) error {
    req := new(models.UserLoginRequest)
    if err := c.Bind(req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request"))
    }

    // Validasi input
    if req.Phone == "" {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Phone required"))
    }

    // Cari user di DB (pseudo code)
    user, err := models.GetUserByPhone(req.Phone)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("User not found"))
    }

    // TODO: Validasi OTP di sini (sementara lewati)

    // Buat JWT token
    token, err := utils.CreateJWTToken(user.ID, user.Role)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create token"))
    }

    return c.JSON(http.StatusOK, echo.Map{
        "message": "Login successful",
        "token":   token,
        "user": echo.Map{
            "id":    user.ID,
            "phone": user.Phone,
            "role":  user.Role,
        },
    })
}
