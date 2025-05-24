package controllers

import (
    "net/http"
    "time"
    "oiya-backend/models"
    "oiya-backend/utils"
    "github.com/labstack/echo/v4"
)

// ChatMessageRequest untuk mengirim pesan chat
type ChatMessageRequest struct {
    SenderID    string `json:"sender_id" validate:"required"`    // ID pengirim (user/admin)
    ReceiverID  string `json:"receiver_id" validate:"required"`  // ID penerima
    Message     string `json:"message" validate:"required"`      // Isi pesan
    ChatType    string `json:"chat_type" validate:"required"`    // "passenger-driver", "driver-admin", "admin-passenger"
}

// SendMessage handler untuk mengirim pesan chat
func SendMessage(c echo.Context) error {
    var req ChatMessageRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid input"))
    }

    chat := models.ChatMessage{
        SenderID:   req.SenderID,
        ReceiverID: req.ReceiverID,
        Message:    req.Message,
        ChatType:   req.ChatType,
        CreatedAt:  time.Now(),
    }

    err := models.SaveChatMessage(&chat)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to save message"))
    }

    return c.JSON(http.StatusOK, utils.NewSuccessResponse("Message sent"))
}

// GetChatMessages handler mengambil percakapan antara dua pihak
func GetChatMessages(c echo.Context) error {
    user1 := c.QueryParam("user1")    // salah satu ID user/admin
    user2 := c.QueryParam("user2")    // ID user/admin lain
    chatType := c.QueryParam("chat_type") // tipe chat

    if user1 == "" || user2 == "" || chatType == "" {
        return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Missing parameters"))
    }

    messages, err := models.GetChatMessagesBetween(user1, user2, chatType)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to get messages"))
    }

    return c.JSON(http.StatusOK, messages)
}
