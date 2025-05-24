package models

import (
    "database/sql"
    "time"
)

// ChatMessage merepresentasikan pesan chat antar pengguna
type ChatMessage struct {
    ID          int64     `json:"id"`
    SenderID    int64     `json:"sender_id"`    // user pengirim
    ReceiverID  int64     `json:"receiver_id"`  // user penerima
    Message     string    `json:"message"`
    CreatedAt   time.Time `json:"created_at"`
    ReadAt      *time.Time `json:"read_at,omitempty"` // nullable, waktu pesan dibaca
}

// SendMessage menyimpan pesan chat baru ke database
func SendMessage(db *sql.DB, msg *ChatMessage) error {
    query := `INSERT INTO chat_messages (sender_id, receiver_id, message, created_at) VALUES (?, ?, ?, NOW())`
    res, err := db.Exec(query, msg.SenderID, msg.ReceiverID, msg.Message)
    if err != nil {
        return err
    }
    msg.ID, err = res.LastInsertId()
    return err
}

// GetChatHistory mengambil daftar pesan antara dua user berdasarkan user ID
func GetChatHistory(db *sql.DB, user1ID, user2ID int64) ([]ChatMessage, error) {
    query := `
        SELECT id, sender_id, receiver_id, message, created_at, read_at
        FROM chat_messages 
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at ASC
    `
    rows, err := db.Query(query, user1ID, user2ID, user2ID, user1ID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var chats []ChatMessage
    for rows.Next() {
        var c ChatMessage
        var readAt sql.NullTime
        err := rows.Scan(&c.ID, &c.SenderID, &c.ReceiverID, &c.Message, &c.CreatedAt, &readAt)
        if err != nil {
            return nil, err
        }
        if readAt.Valid {
            c.ReadAt = &readAt.Time
        } else {
            c.ReadAt = nil
        }
        chats = append(chats, c)
    }
    return chats, nil
}

// MarkMessagesRead tandai semua pesan dari sender ke receiver sebagai sudah dibaca
func MarkMessagesRead(db *sql.DB, senderID, receiverID int64) error {
    query := `UPDATE chat_messages SET read_at = NOW() WHERE sender_id = ? AND receiver_id = ? AND read_at IS NULL`
    _, err := db.Exec(query, senderID, receiverID)
    return err
}
