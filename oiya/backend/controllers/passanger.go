package controllers

import (
    "net/http"
    "oiya-backend/models"
    "oiya-backend/utils"
    "time"

    "github.com/labstack/echo/v4"
)

// RequestPesanDriverRequest mewakili data input pemesanan driver oleh penumpang
type RequestPesanDriverRequest struct {
    PickupLocation   string `json:"pickup_location" validate:"required"`
    Destination      string `json:"destination" validate:"required"`
    Notes            string `json:"notes"`
}

// PesanDriver handler untuk penumpang pesan driver
func PesanDriver(c echo.Context) error {
    var req RequestPesanDriverRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    // Validasi sederhana
    if req.PickupLocation == "" || req.Destination == "" {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Pickup and destination required"))
    }

    // Ambil user ID dari JWT token (middleware)
    userID := utils.GetUserIDFromToken(c)

    // Buat order trip baru (pseudo DB save)
    trip := models.Trip{
        PassengerID:     userID,
        PickupLocation:  req.PickupLocation,
        Destination:     req.Destination,
        Status:          "pending",
        CreatedAt:       time.Now(),
    }

    if err := models.CreateTrip(&trip); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create trip"))
    }

    return c.JSON(http.StatusCreated, utils.NewSuccessResponse("Trip order created"))
}

// TrackingPesanan menampilkan lokasi driver dan info trip aktif
func TrackingPesanan(c echo.Context) error {
    // Ambil user ID
    userID := utils.GetUserIDFromToken(c)

    // Ambil trip aktif user (pseudo DB)
    trip, err := models.GetActiveTripByPassenger(userID)
    if err != nil {
        return c.JSON(http.StatusNotFound, utils.NewErrorResponse("No active trip found"))
    }

    return c.JSON(http.StatusOK, echo.Map{
        "trip_id": trip.ID,
        "driver_location": trip.DriverLocation, // asumsi ada kolom
        "status": trip.Status,
    })
}

// KonfirmasiPembayaran mengkonfirmasi pembayaran (tunai/QRIS)
type KonfirmasiPembayaranRequest struct {
    TripID    int64  `json:"trip_id" validate:"required"`
    Method    string `json:"method" validate:"required"` // "tunai" atau "qris"
    Amount    float64 `json:"amount" validate:"required"`
}

func KonfirmasiPembayaran(c echo.Context) error {
    var req KonfirmasiPembayaranRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    // Proses konfirmasi bayar di DB (pseudo)
    err := models.ConfirmPayment(req.TripID, req.Method, req.Amount)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to confirm payment"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Payment confirmed"))
}

// RiwayatTrip menampilkan daftar trip penumpang
func RiwayatTrip(c echo.Context) error {
    userID := utils.GetUserIDFromToken(c)
    trips, err := models.GetTripsByPassenger(userID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get trips"))
    }

    return c.JSON(http.StatusOK, trips)
}

// BeriRatingRequest input rating untuk driver
type BeriRatingRequest struct {
    TripID  int64   `json:"trip_id" validate:"required"`
    Rating  float32 `json:"rating" validate:"required,min=1,max=5"`
    Comment string  `json:"comment"`
}

func BeriRating(c echo.Context) error {
    var req BeriRatingRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    userID := utils.GetUserIDFromToken(c)

    err := models.AddRating(req.TripID, userID, req.Rating, req.Comment)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to add rating"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Rating submitted"))
}
