package middlewares

import (
    "net/http"
    "strings"
    "oiya-backend/utils"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

// JWTMiddleware fungsi middleware JWT dengan cek role yang diperbolehkan
func JWTMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
    return middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey:    []byte(utils.JwtSecret),
        TokenLookup:   "header:Authorization",
        AuthScheme:    "Bearer",
        ContextKey:    "user",
        ErrorHandlerWithContext: func(err error, c echo.Context) error {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                "error": "Unauthorized: " + err.Error(),
            })
        },
    })
}

// RoleMiddleware untuk cek role user sesuai allowedRoles
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user")
            if user == nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
            }

            token := user.(*middleware.JWTToken)

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
            }

            role, ok := claims["role"].(string)
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Role not found in token"})
            }

            for _, allowed := range allowedRoles {
                if role == allowed {
                    return next(c)
                }
            }

            return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: insufficient permissions"})
        }
    }
}
