package models

import (
    "database/sql"
    "errors"
    "oiya-backend/utils"
    "golang.org/x/crypto/bcrypt"
)

// User struct model user umum untuk passenger, driver, admin
type User struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Phone     string `json:"phone"`
    Email     string `json:"email"`
    Password  string `json:"-"` // tidak ditampilkan di response
    Role      string `json:"role"` // "passenger", "driver", "admin"
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

// CreateUser menyimpan user baru ke database
func CreateUser(db *sql.DB, user *User) error {
    // Hash password sebelum simpan
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := `INSERT INTO users (name, phone, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
    res, err := db.Exec(query, user.Name, user.Phone, user.Email, string(hashedPassword), user.Role)
    if err != nil {
        return err
    }

    user.ID, err = res.LastInsertId()
    return err
}

// GetUserByPhone mengambil user berdasarkan nomor telepon
func GetUserByPhone(db *sql.DB, phone string) (*User, error) {
    user := &User{}
    query := `SELECT id, name, phone, email, password, role, created_at, updated_at FROM users WHERE phone = ?`
    row := db.QueryRow(query, phone)
    err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}

// AuthenticateUser cek login: cocokkan password
func AuthenticateUser(db *sql.DB, phone, password string) (*User, error) {
    user, err := GetUserByPhone(db, phone)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("invalid password")
    }
    return user, nil
}

// UpdateUserPassword update password user
func UpdateUserPassword(db *sql.DB, userID int64, newPassword string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := `UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?`
    _, err = db.Exec(query, string(hashedPassword), userID)
    return err
}
