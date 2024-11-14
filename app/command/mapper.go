package command

import (
	"fmt"
	"github.com/mymmrac/telego"
	"main/db"
	"main/utils"
)

type ResponseMapper struct {
	db   *db.Database
	rate *utils.RateService
}

func NewResponseMapper(db *db.Database, rate *utils.RateService) *ResponseMapper {
	return &ResponseMapper{db: db, rate: rate}
}

func (m *ResponseMapper) MapperCommand(message *telego.Message) string {
	switch message.Text {
	case "/balance", "/balances":
		return m.Balances()
	default:
		if string(message.Text[0]) == "/" {
			return "command or text not found"
		}
		text, err := m.Save(message)
		if err != nil {
			return fmt.Sprintf("ошибка при сохранении информации: %v", err)
		}

		return text
	}
}
