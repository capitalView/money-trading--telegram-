package payload

type NewPayloadMessage struct {
	MessageID    int    `json:"message_id"`
	Date         int64  `json:"date"`
	IsBot        bool   `json:"is_bot"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	Text         string `json:"text"`
}
