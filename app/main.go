package main

import (
	"fmt"
	"main/telegram"
	"main/utils"
	"strconv"
	"time"
)

func main() {
	chatId, _ := strconv.ParseInt(utils.ChatIdAmin, 10, 64)

	done := make(chan bool)
	bot := telegram.NewBotService()

	rate, err := utils.NewRateService()
	bot.SendMessage(chatId, "Bot init")
	if err != nil {
		bot.SendMessage(chatId, fmt.Sprintf("failed to make GET request: %v", err))
		return
	}
	db := utils.NewDatabaseService()

	ticker := time.NewTicker(4 * time.Hour)

	defer func() {
		bot.SendMessage(chatId, "error")
		ticker.Stop()
		bot.StopPulling()
		db.Close()
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				rate.UpdateRates()
				bot.SendMessage(chatId, "Rates updated")
			}
		}
	}()

	updates := bot.StartLongPolling()

	for update := range updates {
		if update.Message != nil {
			text := update.Message.Text
			id := update.Message.Chat.ID
			if text == "/balance" {
				bot.SendMessage(id, db.GetMoney(rate))
				continue
			}
			bot.SendMessage(id, db.SaveInfo(text, rate))
		}
	}

}
