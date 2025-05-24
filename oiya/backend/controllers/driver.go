package controllers

import (
    "net/http"
    "oiya-backend/models"
    "oiya-backend/utils"
    "time"

    "github.com/labstack/echo/v4"
)

// RegisterDriverRequest struct pendaftaran driver lengkap
type RegisterDriverRequest struct {
    Name         string `json:"name" validate:"required"`
    Phone        string `json:"phone" validate:"required"`
    Password     string `json:"password" validate:"required"`
    VehicleType  string `json:"vehicle_type" validate:"required"`
    LicensePlate string `json:"license_plate" validate:"required"`
    Address      string `json:"address"`
    // Tambah field lain sesuai kebutuhan
}

// RegisterDriver handler pendaftaran driver
func RegisterDriver(c echo.Context) error {
    var req RegisterDriverRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    // Cek sudah ada driver dengan nomor hp ini
    exists, err := models.CheckUserExists(req.Phone)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("DB error"))
    }
    if exists {
        return c.JSON(http.StatusConflict, utils.NewErrorResponse("Phone number already registered"))
    }

    // Hash password
    hashedPass, err := utils.HashPassword(req.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to hash password"))
    }

    driver := models.User{
        Name:         req.Name,
        Phone:        req.Phone,
        Password:     hashedPass,
        Role:         "driver",
        VehicleType:  req.VehicleType,
        LicensePlate: req.LicensePlate,
        Address:      req.Address,
        CreatedAt:    time.Now(),
    }

    if err := models.CreateUser(&driver); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create driver"))
    }

    return c.JSON(http.StatusCreated, utils.NewSuccessResponse("Driver registered successfully"))
}

// LoginDriverRequest untuk login driver (hp + OTP)
type LoginDriverRequest struct {
    Phone string `json:"phone" validate:"required"`
    // OTP string `json:"otp" validate:"required"` // Bisa ditambah nanti
}

// LoginDriver handler login driver
func LoginDriver(c echo.Context) error {
    var req LoginDriverRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    user, err := models.GetUserByPhone(req.Phone)
    if err != nil || user.Role != "driver" {
        return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Driver not found"))
    }

    // TODO: Validasi OTP di sini (sementara skip)

    token, err := utils.CreateJWTToken(user.ID, user.Role)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create token"))
    }

    return c.JSON(http.StatusOK, echo.Map{
        "message": "Login successful",
        "token":   token,
        "user": echo.Map{
            "id":           user.ID,
            "name":         user.Name,
            "phone":        user.Phone,
            "vehicle_type": user.VehicleType,
            "license_plate": user.LicensePlate,
        },
    })
}

// DashboardDriver menampilkan data utama dashboard driver
func DashboardDriver(c echo.Context) error {
    driverID := utils.GetUserIDFromToken(c)

    // Ambil status order aktif, saldo, dsb (pseudo DB)
    dashboard, err := models.GetDriverDashboard(driverID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get dashboard data"))
    }

    return c.JSON(http.StatusOK, dashboard)
}

// TerimaOrderRequest untuk terima atau tolak order
type TerimaOrderRequest struct {
    TripID int64  `json:"trip_id" validate:"required"`
    Action string `json:"action" validate:"required"` // "accept" atau "reject"
}

// TerimaOrder handler driver terima/tolak order
func TerimaOrder(c echo.Context) error {
    var req TerimaOrderRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    driverID := utils.GetUserIDFromToken(c)

    if req.Action == "accept" {
        err := models.AcceptOrder(req.TripID, driverID)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to accept order"))
        }
        return c.JSON(http.StatusOK, utils.NewSuccessResponse("Order accepted"))
    } else if req.Action == "reject" {
        err := models.RejectOrder(req.TripID, driverID)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to reject order"))
        }
        return c.JSON(http.StatusOK, utils.NewSuccessResponse("Order rejected"))
    }

    return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid action"))
}

// TrackingPengantaran menampilkan tracking perjalanan antar penumpang
func TrackingPengantaran(c echo.Context) error {
    driverID := utils.GetUserIDFromToken(c)

    trip, err := models.GetActiveTripByDriver(driverID)
    if err != nil {
        return c.JSON(http.StatusNotFound, utils.NewErrorResponse("No active trip found"))
    }

    return c.JSON(http.StatusOK, trip)
}

// RiwayatTripDriver menampilkan riwayat perjalanan driver
func RiwayatTripDriver(c echo.Context) error {
    driverID := utils.GetUserIDFromToken(c)

    trips, err := models.GetTripsByDriver(driverID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get trips"))
    }

    return c.JSON(http.StatusOK, trips)
}

// TopUpSaldoRequest untuk isi saldo membership driver
type TopUpSaldoRequest struct {
    Amount float64 `json:"amount" validate:"required"`
}

func TopUpSaldo(c echo.Context) error {
    var req TopUpSaldoRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    driverID := utils.GetUserIDFromToken(c)

    err := models.TopUpDriverBalance(driverID, req.Amount)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to top up saldo"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Top up successful"))
}
