package middleware

import (
    "net/http"
    "os"

    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
    return middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey:  []byte(os.Getenv("JWT_SECRET")),
        TokenLookup: "header:Authorization",
        AuthScheme:  "Bearer",
        ErrorHandlerWithContext: func(err error, c echo.Context) error {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
        },
    })
}

func GetUserID(c echo.Context) int {
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)
    return int(claims["user_id"].(float64))
}
