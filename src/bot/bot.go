package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"strconv"
	"strings"
	"../config"
)

var asked bool

func ListenForUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel)  {
	for update := range updates {

		fmt.Println(update.Message)

		if update.Message != nil {
			cmd := update.Message.Command()

			if cmd == "" {
				if asked {
					channelId, err := strconv.ParseInt(config.Toml.Bot.ChannelId, 10, 64);

					if err != nil {
						fmt.Println(err)
					}

					msg := tgbotapi.NewForward(channelId, update.Message.Chat.ID, update.Message.MessageID)

					bot.Send(msg)

					confirmMsg := newMessage(
						update.Message.Chat.ID,
						"<b>Ваш вопрос принят!</b>",
						"html")

					bot.Send(confirmMsg)

					asked = false
				}
			} else {
				switch cmd {
				case "start", "menu":
					msg := newMessage(update.Message.Chat.ID, "<b>Меню:</b>", "html")

					keyboard := tgbotapi.InlineKeyboardMarkup{}
					menu := getMenu()

					for _, item := range menu {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(item.Name, item.Alias)
						row = append(row, btn)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}

					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				}
			}
		} else {
			if update.CallbackQuery != nil {

				if (strings.Contains(update.CallbackQuery.Data, "faq_")) {
					s := strings.Split(update.CallbackQuery.Data, "_")
					id, err := strconv.Atoi(s[1])

					if err != nil {
						fmt.Println(err)
					}

					question := getQuestion(id)

					msg := newMessage(
						update.CallbackQuery.Message.Chat.ID,
						"<b>" + question.Question + "</b>" + "\n\n" + question.Answer,
						"html")

					bot.Send(msg)
				}

				switch update.CallbackQuery.Data {
				case "schedule":
					schedule := getSchedule()

					msg := newMessage(
						update.CallbackQuery.Message.Chat.ID,
						"<b>Расписание:</b>" + "\n\n" + schedule.Value,
						"html")

					bot.Send(msg)
				case "ask":
					msg := newMessage(
						update.CallbackQuery.Message.Chat.ID,
						"<b>У вас есть вопрос? Задайте его спикерам</b>",
						"html")

					bot.Send(msg)

					asked = true
				case "faq":
					msg := newMessage(update.CallbackQuery.Message.Chat.ID, "<b>Часто задаваемые вопросы:</b>", "html")

					keyboard := tgbotapi.InlineKeyboardMarkup{}
					questions := getFaq()


					for _, item := range questions {
						var row []tgbotapi.InlineKeyboardButton
						btn := tgbotapi.NewInlineKeyboardButtonData(item.Question, "faq_" + strconv.Itoa(item.Id))
						row = append(row, btn)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}

					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				}
			}
		}
	}
}

func newMessage(chatId int64, text string, parseMode string) tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatId,
			ReplyToMessageID: 0,
		},
		Text: text,
		ParseMode: parseMode,
		DisableWebPagePreview: false,
	}
}