package panel

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"../helper"
)

type BotPanel struct {
	BotApi *tgbotapi.BotAPI
	Update tgbotapi.Update
}

func (bot *BotPanel) ListenPanelCmds() {
	switch bot.Update.Message.Command() {
	case "panel_invite":
	case "panel_help":
		bot.panelHelp()
	case "panel_import_test":
		bot.importTests()
	}
}

func (bot *BotPanel) ListenPanelMsgs() {

}

func (bot *BotPanel) inviteUser() {

}

func (bot *BotPanel) panelHelp() {
	fmt.Println()
	fmt.Println("help")
	fmt.Println()
}

func (bot *BotPanel) importTests() {
	msg := helper.NewMessage(bot.Update.Message.Chat.ID, "Загрузите файл с тестами в формате xml", "html")
	bot.BotApi.Send(msg)
}