package payload

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	"main/db"
)

type Payload struct {
	*db.Database
	payload NewPayloadMessage
}

func NewPayload(db *db.Database, message *telego.Message) *Payload {
	payload := NewPayloadMessage{
		MessageID:    message.MessageID,
		Date:         message.Date,
		IsBot:        message.From.IsBot,
		Username:     message.From.Username,
		LanguageCode: message.From.LanguageCode,
		Text:         message.Text,
	}
	return &Payload{Database: db, payload: payload}
}

func (data *Payload) GetTransactionId(ctx context.Context, messageId int) (int, error) {
	rows, err := data.Query(ctx, GetTransactionId, messageId)
	if err != nil {
		fmt.Println(err)
	}

	var trId int
	for rows.Next() {
		if err := rows.Scan(&trId); err != nil {
			return 0, fmt.Errorf("Ошибка при чтении строки: %v\n", err)
		}
	}

	return trId, nil
}

func (data *Payload) SavePayload(ctx context.Context, transactionID int) error {
	return data.Execute(ctx, InsertPayload, transactionID, data.payload)
}
func (data *Payload) UpdatePayload(ctx context.Context, transactionID int) {
	data.Execute(ctx, UpdatePayload, transactionID, data.payload)
}
