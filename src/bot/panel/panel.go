package panel

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"../helper"
	//"github.com/tealeg/xlsx"
	"../../config"
	"io/ioutil"
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
	case "panel_import_test":
		bot.importTests()
	}
}

func (bot *BotPanel) ListenPanelMsgs() {
	fmt.Println("sdfsfdsf")

	if _, ok := MsgCmd[bot.Update.Message.From.ID]; ok && MsgCmd[bot.Update.Message.From.ID].Status {
		switch MsgCmd[bot.Update.Message.From.ID].Cmd {
		case "panel_import_test":
			bot.readImportedFile()
		}
	}
}

func (bot *BotPanel) inviteUser() {

}

func (bot *BotPanel) panelHelp() {

}

func (bot *BotPanel) importTests() {
	msg := helper.NewMessage(bot.Update.Message.Chat.ID, "Загрузите файл с тестами в формате xml", "html")
	bot.BotApi.Send(msg)
}

func (bot *BotPanel) readImportedFile() {
	fileId := bot.Update.Message.Document.FileID
	file, err := bot.BotApi.GetFile(tgbotapi.FileConfig{fileId})

	if err == nil {
		excelFileName := "https://api.telegram.org/file/bot" + config.Toml.Bot.Token + "/" + file.FilePath
		fmt.Println(excelFileName)

		//xlFile, err := xlsx.OpenFile(excelFileName)

		bytes, err := ioutil.ReadFile(excelFileName)

		fmt.Println()
		fmt.Println(bytes)
		fmt.Println(err)
		fmt.Println()

		//if err != nil {
		//
		//}
		//for _, sheet := range xlFile.Sheets {
		//	for _, row := range sheet.Rows {
		//		for _, cell := range row.Cells {
		//			text := cell.String()
		//			fmt.Printf("%s\n", text)
		//		}
		//	}
		//}
	}
}

func (bot *BotPanel) IsMessaging() bool {
	if _, ok := MsgCmd[bot.Update.Message.From.ID]; ok && MsgCmd[bot.Update.Message.From.ID].Status {
		return true
	}

	return false
}