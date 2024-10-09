package telegram

import (
	"fmt"
	"github.com/mymmrac/telego"
	"main/utils"
	"os"
)

type BotService struct {
	bot *telego.Bot
}

func NewBotService() *BotService {
	bot, err := telego.NewBot(utils.TokenTelegram, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &BotService{bot: bot}
}

func (bs *BotService) SendMessage(id int64, text string) {
	bs.bot.SendMessage(&telego.SendMessageParams{
		ChatID: telego.ChatID{
			ID: id,
		},
		Text: text,
	})
}

func (bs *BotService) StartLongPolling() <-chan telego.Update {
	updates, _ := bs.bot.UpdatesViaLongPolling(&telego.GetUpdatesParams{
		Offset: -1,
	})

	return updates

}
func (bs *BotService) StopPulling() {
	bs.bot.StopLongPolling()
}
