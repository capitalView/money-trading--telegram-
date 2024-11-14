package main

import (
	"fmt"
	"main/command"
	"main/db"
	"main/telegram"
	"main/utils"
	"strconv"
	"time"
)

var ChatIdAmin = utils.GoDotEnvVariable("CHAT_ID")

func main() {
	chatId, _ := strconv.ParseInt(ChatIdAmin, 10, 64)

	done := make(chan bool)
	bot := telegram.NewBotService()

	rate, err := utils.NewRateService()
	bot.Send(chatId, "Bot init")
	if err != nil {
		bot.Send(chatId, fmt.Sprintf("failed to make GET request: %v", err))
		return
	}

	ticker := time.NewTicker(4 * time.Hour)

	database := db.NewDatabaseService()
	defer func() {
		ticker.Stop()
		bot.StopPulling()
		database.Close()
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				rate.UpdateRates()
				bot.Send(chatId, "Rates updated")
			}
		}
	}()

	mapper := command.NewResponseMapper(database, rate)

	updates := bot.StartLongPolling()

	for update := range updates {
		if update.EditedMessage != nil {
			id := update.EditedMessage.Chat.ID
			bot.Send(id, mapper.Edit(update.EditedMessage))
			bot.Delete(id, update.EditedMessage.MessageID+1)
		}
		if update.Message != nil {
			id := update.Message.Chat.ID
			bot.Send(id, mapper.MapperCommand(update.Message))
		}
	}
}
