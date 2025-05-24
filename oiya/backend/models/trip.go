package models

import (
    "database/sql"
    "time"
)

// Trip merepresentasikan data perjalanan ojek
type Trip struct {
    ID            int64     `json:"id"`
    PassengerID   int64     `json:"passenger_id"`
    DriverID      int64     `json:"driver_id"`
    StartLocation string    `json:"start_location"` // bisa JSON string lat-long atau alamat
    EndLocation   string    `json:"end_location"`
    StartTime     time.Time `json:"start_time"`
    EndTime       time.Time `json:"end_time"`
    Status        string    `json:"status"`      // pending, ongoing, completed, cancelled
    Fare          float64   `json:"fare"`        // tarif perjalanan
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// CreateTrip menyimpan data trip baru ke database
func CreateTrip(db *sql.DB, trip *Trip) error {
    query := `INSERT INTO trips (passenger_id, driver_id, start_location, end_location, start_time, end_time, status, fare, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
    res, err := db.Exec(query, trip.PassengerID, trip.DriverID, trip.StartLocation, trip.EndLocation, trip.StartTime, trip.EndTime, trip.Status, trip.Fare)
    if err != nil {
        return err
    }

    trip.ID, err = res.LastInsertId()
    return err
}

// GetTripByID mengambil trip berdasarkan ID
func GetTripByID(db *sql.DB, id int64) (*Trip, error) {
    trip := &Trip{}
    query := `SELECT id, passenger_id, driver_id, start_location, end_location, start_time, end_time, status, fare, created_at, updated_at FROM trips WHERE id = ?`
    row := db.QueryRow(query, id)
    err := row.Scan(&trip.ID, &trip.PassengerID, &trip.DriverID, &trip.StartLocation, &trip.EndLocation, &trip.StartTime, &trip.EndTime, &trip.Status, &trip.Fare, &trip.CreatedAt, &trip.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return trip, nil
}

// GetTripsByPassenger mengambil daftar trip berdasarkan penumpang
func GetTripsByPassenger(db *sql.DB, passengerID int64) ([]Trip, error) {
    query := `SELECT id, passenger_id, driver_id, start_location, end_location, start_time, end_time, status, fare, created_at, updated_at 
              FROM trips WHERE passenger_id = ? ORDER BY start_time DESC`
    rows, err := db.Query(query, passengerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var trips []Trip
    for rows.Next() {
        var trip Trip
        err := rows.Scan(&trip.ID, &trip.PassengerID, &trip.DriverID, &trip.StartLocation, &trip.EndLocation, &trip.StartTime, &trip.EndTime, &trip.Status, &trip.Fare, &trip.CreatedAt, &trip.UpdatedAt)
        if err != nil {
            return nil, err
        }
        trips = append(trips, trip)
    }
    return trips, nil
}

// GetTripsByDriver mengambil daftar trip berdasarkan driver
func GetTripsByDriver(db *sql.DB, driverID int64) ([]Trip, error) {
    query := `SELECT id, passenger_id, driver_id, start_location, end_location, start_time, end_time, status, fare, created_at, updated_at 
              FROM trips WHERE driver_id = ? ORDER BY start_time DESC`
    rows, err := db.Query(query, driverID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var trips []Trip
    for rows.Next() {
        var trip Trip
        err := rows.Scan(&trip.ID, &trip.PassengerID, &trip.DriverID, &trip.StartLocation, &trip.EndLocation, &trip.StartTime, &trip.EndTime, &trip.Status, &trip.Fare, &trip.CreatedAt, &trip.UpdatedAt)
        if err != nil {
            return nil, err
        }
        trips = append(trips, trip)
    }
    return trips, nil
}

// UpdateTripStatus update status trip dan waktu akhir jika diperlukan
func UpdateTripStatus(db *sql.DB, tripID int64, status string, endTime *time.Time) error {
    var query string
    if endTime != nil {
        query = `UPDATE trips SET status = ?, end_time = ?, updated_at = NOW() WHERE id = ?`
        _, err := db.Exec(query, status, *endTime, tripID)
        return err
    } else {
        query = `UPDATE trips SET status = ?, updated_at = NOW() WHERE id = ?`
        _, err := db.Exec(query, status, tripID)
        return err
    }
}
