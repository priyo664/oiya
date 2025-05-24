package models

import (
    "database/sql"
    "time"
)

// Payment merepresentasikan data pembayaran untuk trip
type Payment struct {
    ID          int64     `json:"id"`
    TripID      int64     `json:"trip_id"`
    Amount      float64   `json:"amount"`
    Method      string    `json:"method"`      // e.g. "cash", "QRIS", "e-wallet"
    Status      string    `json:"status"`      // e.g. "pending", "completed", "failed"
    PaidAt      time.Time `json:"paid_at"`     // waktu pembayaran, bisa nol jika belum bayar
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePayment menyimpan data pembayaran baru
func CreatePayment(db *sql.DB, p *Payment) error {
    query := `INSERT INTO payments (trip_id, amount, method, status, paid_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
    res, err := db.Exec(query, p.TripID, p.Amount, p.Method, p.Status, p.PaidAt)
    if err != nil {
        return err
    }

    p.ID, err = res.LastInsertId()
    return err
}

// GetPaymentByTripID mengambil data pembayaran berdasarkan trip_id
func GetPaymentByTripID(db *sql.DB, tripID int64) (*Payment, error) {
    payment := &Payment{}
    query := `SELECT id, trip_id, amount, method, status, paid_at, created_at, updated_at FROM payments WHERE trip_id = ?`
    row := db.QueryRow(query, tripID)
    err := row.Scan(&payment.ID, &payment.TripID, &payment.Amount, &payment.Method, &payment.Status, &payment.PaidAt, &payment.CreatedAt, &payment.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return payment, nil
}

// UpdatePaymentStatus update status pembayaran dan waktu bayar jika perlu
func UpdatePaymentStatus(db *sql.DB, paymentID int64, status string, paidAt *time.Time) error {
    var query string
    if paidAt != nil {
        query = `UPDATE payments SET status = ?, paid_at = ?, updated_at = NOW() WHERE id = ?`
        _, err := db.Exec(query, status, *paidAt, paymentID)
        return err
    } else {
        query = `UPDATE payments SET status = ?, updated_at = NOW() WHERE id = ?`
        _, err := db.Exec(query, status, paymentID)
        return err
    }
}
