package utils

import (
	"fmt"
	"github.com/mymmrac/telego"
)

type ResponseMapper struct {
	db   *DatabaseService
	rate *RateService
}

func NewResponseMapperServices(db *DatabaseService, rate *RateService) *ResponseMapper {
	return &ResponseMapper{db: db, rate: rate}
}

func (r *ResponseMapper) save(message *telego.Message) (string, error) {
	messagePayload := GetMessagePayload(message)
	id, err := r.db.SaveInfo(messagePayload.Text, r.rate)
	if err != nil {
		return "", err
	}
	r.db.SavePayload(messagePayload, id)
	return r.db.GetAll(r.rate), nil
}

func (r *ResponseMapper) Update(message *telego.Message) string {
	messagePayload := GetMessagePayload(message)
	id, err := r.db.GetTransactionId(messagePayload.MessageID)
	fmt.Println(id)
	if err != nil {
		return fmt.Sprintf("%x", err)
	}
	result, err := ParseText(messagePayload.Text)
	if err != nil {
		return fmt.Sprintf("%x", err)
	}
	errorUpdate := r.db.UpdateTransaction(result, id)
	if errorUpdate != nil {
		return fmt.Sprintf("%x", errorUpdate)
	}
	r.db.UpdatePayload(messagePayload, id)
	return r.db.GetAll(r.rate)
}

func (r *ResponseMapper) MapperCommand(message *telego.Message) string {
	messageText := message.Text
	switch messageText {
	case "/balance", "/balances":
		return r.db.GetAll(r.rate)
	default:
		if string(messageText[0]) == "/" {
			return "command or text not found"
		}
		text, err := r.save(message)
		if err != nil {
			return fmt.Sprintf("ошибка при сохранении информации: %v", err)
		}

		return text
	}
}
