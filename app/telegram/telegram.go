package telegram

import (
	"fmt"
	"github.com/mymmrac/telego"
	"main/utils"
	"os"
)

type BotService struct {
	*telego.Bot
}

var TokenTelegram = utils.GoDotEnvVariable("TELEGRAM_TOKEN")

func NewBotService() *BotService {
	bot, err := telego.NewBot(TokenTelegram, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &BotService{bot}
}

func (b *BotService) Send(id int64, text string) {
	b.SendMessage(&telego.SendMessageParams{
		ChatID: telego.ChatID{
			ID: id,
		},
		Text:      text,
		ParseMode: "Markdown",
	})
}

func (b *BotService) Delete(id int64, messageId int) {
	b.DeleteMessage(&telego.DeleteMessageParams{
		ChatID: telego.ChatID{
			ID: id,
		},
		MessageID: messageId,
	})
}

func (b *BotService) StartLongPolling() <-chan telego.Update {
	updates, _ := b.UpdatesViaLongPolling(&telego.GetUpdatesParams{
		Offset: -1,
	})

	return updates

}
func (b *BotService) StopPulling() {
	b.StopLongPolling()
}
