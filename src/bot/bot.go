package bot

import (
	"../config"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"strconv"
)

func ListenForUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel)  {
	for update := range updates {
		if update.Message != nil {
			cmd := update.Message.Command()

			if cmd == "" {
				//для обычных сообщений
			} else {
				switch cmd {
				case "start":
					msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, config.Toml.Bot.StartPic)

					keyboard := tgbotapi.InlineKeyboardMarkup{}
					var row []tgbotapi.InlineKeyboardButton
					btn := tgbotapi.NewInlineKeyboardButtonData("Play", "/play")
					row = append(row, btn)
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				case "play":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для начала выберите кем вы являетесь")

					keyboard := tgbotapi.InlineKeyboardMarkup{}
					categories := getCategories()

					for _, category := range categories {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(category.Name, strconv.Itoa(category.Id))
						row = append(row, btn)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}

					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				}
			}
		} else {
			if update.CallbackQuery != nil {
				fmt.Println("here")
				fmt.Println(update.CallbackQuery)
				fmt.Println("here")

				//category := update.CallbackQuery.Data
				//categoriesMap[update.CallbackQuery.From.ID] = category
				//bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отлично, я запомнил"))
			}
		}
	}
}