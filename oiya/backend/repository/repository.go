package repository

import (
    "database/sql"
    "errors"
    "oiya/models"
    "time"
)

// UserRepository interface untuk operasi data user
type UserRepository interface {
    CreateUser(user *models.User) error
    GetUserByID(id int64) (*models.User, error)
    GetUserByPhone(phone string) (*models.User, error)
    UpdateUser(user *models.User) error
}

// TripRepository interface untuk operasi data trip
type TripRepository interface {
    CreateTrip(trip *models.Trip) error
    GetTripByID(id int64) (*models.Trip, error)
    UpdateTrip(trip *models.Trip) error
    ListTripsByUser(userID int64) ([]models.Trip, error)
}

// PaymentRepository interface untuk operasi data pembayaran
type PaymentRepository interface {
    CreatePayment(payment *models.Payment) error
    GetPaymentByTripID(tripID int64) (*models.Payment, error)
    UpdatePaymentStatus(paymentID int64, status string, paidAt *time.Time) error
}

// ChatRepository interface untuk operasi chat
type ChatRepository interface {
    SendMessage(msg *models.ChatMessage) error
    GetChatHistory(user1ID, user2ID int64) ([]models.ChatMessage, error)
    MarkMessagesRead(senderID, receiverID int64) error
}

// --- Implementasi concret struct yang memakai *sql.DB ---

type userRepo struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *models.User) error {
    query := `INSERT INTO users (name, phone, email, password_hash, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
    res, err := r.db.Exec(query, user.Name, user.Phone, user.Email, user.PasswordHash, user.Role)
    if err != nil {
        return err
    }
    user.ID, err = res.LastInsertId()
    return err
}

func (r *userRepo) GetUserByID(id int64) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, name, phone, email, password_hash, role, created_at, updated_at FROM users WHERE id = ?`
    row := r.db.QueryRow(query, id)
    err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}

func (r *userRepo) GetUserByPhone(phone string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, name, phone, email, password_hash, role, created_at, updated_at FROM users WHERE phone = ?`
    row := r.db.QueryRow(query, phone)
    err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}

func (r *userRepo) UpdateUser(user *models.User) error {
    query := `UPDATE users SET name = ?, phone = ?, email = ?, password_hash = ?, role = ?, updated_at = NOW() WHERE id = ?`
    _, err := r.db.Exec(query, user.Name, user.Phone, user.Email, user.PasswordHash, user.Role, user.ID)
    return err
}

// Implementasi repository lain (TripRepository, PaymentRepository, ChatRepository) dibuat serupa
// Kalau kamu mau, aku bisa buatkan juga implementasi lengkap untuk repository lain.

