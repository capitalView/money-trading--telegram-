package main

import (
	"encoding/json"
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

	mapper := utils.NewResponseMapperServices(db, rate)

	updates := bot.StartLongPolling()

	for update := range updates {
		if update.EditedMessage != nil {
			id := update.EditedMessage.Chat.ID
			print(id, update.EditedMessage)
			bot.SendMessage(id, mapper.Update(update.EditedMessage))
		}
		if update.Message != nil {
			id := update.Message.Chat.ID
			bot.SendMessage(id, mapper.MapperCommand(update.Message))
		}
	}
}

func print(data ...any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
