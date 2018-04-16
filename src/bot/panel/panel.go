package panel

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var MsgCmd = map[int]*MsgTracker{}

type MsgTracker struct {
	Cmd    string
	Status bool
}

type BotPanel struct {
	BotApi *tgbotapi.BotAPI
	Update tgbotapi.Update
}

func (bot *BotPanel) ListenPanelCmds() {
	MsgCmd[bot.Update.Message.From.ID] = &MsgTracker{bot.Update.Message.Command(), true}

	switch bot.Update.Message.Command() {
	case "panel_invite":
	case "panel_help":
		bot.panelHelp()
	}
}

func (bot *BotPanel) ListenPanelMsgs() {
	if _, ok := MsgCmd[bot.Update.Message.From.ID]; ok && MsgCmd[bot.Update.Message.From.ID].Status {
		switch MsgCmd[bot.Update.Message.From.ID].Cmd {
		case "":
		}
	}
}

func (bot *BotPanel) inviteUser() {

}

func (bot *BotPanel) panelHelp() {

}

func (bot *BotPanel) IsMessaging() bool {
	if _, ok := MsgCmd[bot.Update.Message.From.ID]; ok && MsgCmd[bot.Update.Message.From.ID].Status {
		return true
	}

	return false
}