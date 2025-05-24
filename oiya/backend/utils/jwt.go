package middleware

import (
    "errors"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your-very-secret-key") // Ganti dengan env var di production

// GenerateToken membuat JWT token dengan payload userID dan role
func GenerateToken(userID int64, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(24 * time.Hour).Unix(), // token berlaku 24 jam
        "iat":     time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// JWTMiddleware middleware untuk validasi JWT token di header Authorization
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
        }
        tokenStr := parts[1]

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        })

        if err != nil {
            return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
        }
        if !token.Valid {
            return echo.NewHTTPError(http.StatusUnauthorized, "token is not valid")
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return echo.NewHTTPError(http.StatusUnauthorized, "cannot read claims")
        }

        // Set user info ke context agar bisa dipakai di handler
        c.Set("user_id", int64(claims["user_id"].(float64)))
        c.Set("role", claims["role"].(string))

        return next(c)
    }
}

// RoleMiddleware middleware untuk validasi role user
func RoleMiddleware(requiredRole string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            role, ok := c.Get("role").(string)
            if !ok || role != requiredRole {
                return echo.NewHTTPError(http.StatusForbidden, "forbidden: insufficient role")
            }
            return next(c)
        }
    }
}
