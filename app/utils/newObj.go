package utils

import "github.com/mymmrac/telego"

type NewPayloadMessage struct {
	MessageID    int    `json:"message_id"`
	Date         int64  `json:"date"`
	IsBot        bool   `json:"is_bot"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	Text         string `json:"text"`
}

func GetMessagePayload(payload *telego.Message) NewPayloadMessage {
	return NewPayloadMessage{
		MessageID:    payload.MessageID,
		Date:         payload.Date,
		IsBot:        payload.From.IsBot,
		Username:     payload.From.Username,
		LanguageCode: payload.From.LanguageCode,
		Text:         payload.Text,
	}
}
