package routes

import (
    "oiya/handlers"
    "oiya/middleware"

    "github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
    // Public routes (no auth)
    e.POST("/api/auth/register", handlers.Register)
    e.POST("/api/auth/login", handlers.Login)
    e.POST("/api/auth/verify-otp", handlers.VerifyOTP)

    // Group: Passenger (require auth + role passenger)
    passengerGroup := e.Group("/api/passenger")
    passengerGroup.Use(middleware.JWTMiddleware)
    passengerGroup.Use(middleware.RoleMiddleware("passenger"))
    {
        passengerGroup.GET("/dashboard", handlers.PassengerDashboard)
        passengerGroup.POST("/order", handlers.CreateOrder)
        passengerGroup.GET("/orders", handlers.ListOrders)
        passengerGroup.GET("/orders/:id", handlers.GetOrder)
        passengerGroup.POST("/orders/:id/rate", handlers.RateDriver)
        passengerGroup.POST("/payment", handlers.MakePayment)
        passengerGroup.GET("/chat/:driverId", handlers.GetChatWithDriver)
        passengerGroup.POST("/chat/send", handlers.SendChatMessage)
    }

    // Group: Driver (require auth + role driver)
    driverGroup := e.Group("/api/driver")
    driverGroup.Use(middleware.JWTMiddleware)
    driverGroup.Use(middleware.RoleMiddleware("driver"))
    {
        driverGroup.GET("/dashboard", handlers.DriverDashboard)
        driverGroup.POST("/accept-order/:orderId", handlers.AcceptOrder)
        driverGroup.POST("/reject-order/:orderId", handlers.RejectOrder)
        driverGroup.GET("/trips", handlers.ListTrips)
        driverGroup.GET("/trips/:id", handlers.GetTrip)
        driverGroup.POST("/topup", handlers.TopUpMembership)
        driverGroup.GET("/chat/:passengerId", handlers.GetChatWithPassenger)
        driverGroup.POST("/chat/send", handlers.SendChatMessage)
    }

    // Group: Admin (require auth + role admin)
    adminGroup := e.Group("/api/admin")
    adminGroup.Use(middleware.JWTMiddleware)
    adminGroup.Use(middleware.RoleMiddleware("admin"))
    {
        adminGroup.GET("/dashboard", handlers.AdminDashboard)
        adminGroup.GET("/users", handlers.ListUsers)
        adminGroup.PUT("/users/:id", handlers.UpdateUser)
        adminGroup.GET("/orders", handlers.ListAllOrders)
        adminGroup.PUT("/orders/:id", handlers.UpdateOrder)
        adminGroup.GET("/payments", handlers.ListPayments)
        adminGroup.GET("/reports", handlers.GetReports)
        adminGroup.PUT("/settings", handlers.UpdateSettings)
        adminGroup.POST("/broadcast", handlers.BroadcastMessage)
        adminGroup.POST("/banner", handlers.AddBanner)
        adminGroup.GET("/chat/:userId", handlers.GetChatWithUser)
        adminGroup.POST("/chat/send", handlers.SendChatMessage)
    }
}
