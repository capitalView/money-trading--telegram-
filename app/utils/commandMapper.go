package utils

import "fmt"

type ResponseMapper struct {
	db   *DatabaseService
	rate *RateService
}

func NewResponseMapperServices(db *DatabaseService, rate *RateService) *ResponseMapper {
	return &ResponseMapper{db: db, rate: rate}
}

func (r *ResponseMapper) MapperCommand(text string) string {
	switch text {
	case "/balance", "/balances":
		return r.db.GetAll(r.rate)
	default:
		if string(text[0]) == "/" {
			return "команда или текст не найден"
		}
		text, err := r.db.SaveInfo(text, r.rate)
		if err != nil {
			return fmt.Sprintf("ошибка при сохранении информации: %v", err)
		}

		return text
	}
}
