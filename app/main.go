package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
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

	database := db.NewDatabaseService(ctx)

	defer func() {
		ticker.Stop()
		bot.StopPulling()
		database.Close(ctx)
		cancel()
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				rate.UpdateRates()
			}
		}
	}()

	mapper := command.NewResponseMapper(database, rate)

	updates := bot.StartLongPolling()

	for update := range updates {
		if update.EditedMessage != nil {
			id := update.EditedMessage.Chat.ID
			bot.Send(id, mapper.Edit(ctx, update.EditedMessage))
			bot.Delete(id, update.EditedMessage.MessageID+1)
		}
		if update.Message != nil {
			id := update.Message.Chat.ID
			bot.Send(id, mapper.MapperCommand(ctx, update.Message))
		}
	}
}
