package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword menerima plain password dan mengembalikan hashed password
func HashPassword(password string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashed), nil
}

// CheckPassword membandingkan plain password dengan hashed password
func CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
