package controllers

import (
    "net/http"
    "oiya-backend/models"
    "oiya-backend/utils"
    "github.com/labstack/echo/v4"
)

// LoginAdminRequest untuk login admin
type LoginAdminRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

// LoginAdmin handler login admin
func LoginAdmin(c echo.Context) error {
    var req LoginAdminRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    admin, err := models.GetAdminByUsername(req.Username)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Admin not found"))
    }

    if !utils.CheckPasswordHash(req.Password, admin.Password) {
        return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Wrong password"))
    }

    token, err := utils.CreateJWTToken(admin.ID, "admin")
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create token"))
    }

    return c.JSON(http.StatusOK, echo.Map{
        "message": "Login successful",
        "token":   token,
        "admin": echo.Map{
            "id":       admin.ID,
            "username": admin.Username,
            "name":     admin.Name,
        },
    })
}

// DashboardAdmin menampilkan ringkasan statistik & grafik
func DashboardAdmin(c echo.Context) error {
    stats, err := models.GetAdminDashboardStats()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get dashboard data"))
    }
    return c.JSON(http.StatusOK, stats)
}

// ManajemenUser menampilkan daftar dan kelola user (penumpang & driver)
func ManajemenUser(c echo.Context) error {
    role := c.QueryParam("role") // optional: "passenger" or "driver"

    users, err := models.GetUsersByRole(role)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get users"))
    }
    return c.JSON(http.StatusOK, users)
}

// UpdateUserRequest untuk update data user
type UpdateUserRequest struct {
    Name         *string `json:"name"`
    Phone        *string `json:"phone"`
    VehicleType  *string `json:"vehicle_type,omitempty"`
    LicensePlate *string `json:"license_plate,omitempty"`
    Status       *string `json:"status,omitempty"` // aktif/nonaktif, dsb
}

// UpdateUser handler update user info
func UpdateUser(c echo.Context) error {
    userID := c.Param("id")
    var req UpdateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    err := models.UpdateUser(userID, req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update user"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("User updated successfully"))
}

// ManajemenOrder menampilkan dan kelola order
func ManajemenOrder(c echo.Context) error {
    orders, err := models.GetAllOrders()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get orders"))
    }
    return c.JSON(http.StatusOK, orders)
}

// UpdateOrderRequest untuk update order status
type UpdateOrderRequest struct {
    Status string `json:"status" validate:"required"` // contoh: pending, accepted, completed, cancelled
}

// UpdateOrder handler update status order
func UpdateOrder(c echo.Context) error {
    orderID := c.Param("id")
    var req UpdateOrderRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    err := models.UpdateOrderStatus(orderID, req.Status)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update order status"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Order status updated"))
}

// ManajemenPembayaran menampilkan transaksi dan status pembayaran
func ManajemenPembayaran(c echo.Context) error {
    payments, err := models.GetAllPayments()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get payments"))
    }
    return c.JSON(http.StatusOK, payments)
}

// LaporanStatistik menampilkan grafik dan laporan perjalanan & pendapatan
func LaporanStatistik(c echo.Context) error {
    report, err := models.GetReportData()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get report data"))
    }
    return c.JSON(http.StatusOK, report)
}

// PengaturanSistemRequest untuk update setting sistem
type PengaturanSistemRequest struct {
    TarifBaseFare     *float64 `json:"tarif_base_fare,omitempty"`
    TarifPerKm        *float64 `json:"tarif_per_km,omitempty"`
    MaxDistance       *float64 `json:"max_distance,omitempty"`
    // dst sesuai kebutuhan
}

// PengaturanSistem handler update pengaturan sistem
func PengaturanSistem(c echo.Context) error {
    var req PengaturanSistemRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    err := models.UpdateSystemSettings(req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update settings"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Settings updated"))
}

// ChatAdminDriver menampilkan chat admin <-> driver (mock example)
func ChatAdminDriver(c echo.Context) error {
    // Implementasi chat backend bisa pakai DB atau real-time server (MQTT/WebSocket)
    chats, err := models.GetChatsAdminDriver()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get chat"))
    }
    return c.JSON(http.StatusOK, chats)
}

// BroadcastMessageRequest untuk kirim pesan broadcast
type BroadcastMessageRequest struct {
    Target   string `json:"target" validate:"required"` // "passenger", "driver", atau "all"
    Message  string `json:"message" validate:"required"`
}

// BroadcastMessage handler kirim pesan broadcast
func BroadcastMessage(c echo.Context) error {
    var req BroadcastMessageRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    err := models.SendBroadcastMessage(req.Target, req.Message)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to send broadcast"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Broadcast sent"))
}

// BannerRequest untuk pasang banner iklan di homepage penumpang
type BannerRequest struct {
    ImageURL string `json:"image_url" validate:"required"`
    LinkURL  string `json:"link_url" validate:"required"`
}

// AddBanner handler tambah banner baru
func AddBanner(c echo.Context) error {
    var req BannerRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    err := models.AddBanner(req.ImageURL, req.LinkURL)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to add banner"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Banner added"))
}

// GetBanners menampilkan list banner aktif
func GetBanners(c echo.Context) error {
    banners, err := models.GetBanners()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get banners"))
    }
    return c.JSON(http.StatusOK, banners)
}
