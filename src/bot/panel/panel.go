package panel

import "github.com/go-telegram-bot-api/telegram-bot-api"

type BotPanel struct {
	BotApi *tgbotapi.BotAPI
	Update tgbotapi.Update
}

func (bot *BotPanel) PanelInit() {
	switch bot.Update.Message.Command() {
	case "panel_help":
		bot.panelHelp()
	case "import_test":
		bot.importTests()
	}
}

func (bot *BotPanel) panelHelp() {

}

func (bot *BotPanel) importTests ()  {

}