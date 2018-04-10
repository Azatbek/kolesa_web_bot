package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"strconv"
	"strings"
	"../config"
	"../db"
	"time"
)

var (
	asked    bool
	variants []string
)

var Chats = map[string]*Quiz{}

type BotApi struct {
	BotApi  *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
	Update  tgbotapi.Update
}

type Quiz struct {
	User                    string
	Index, Score, LastMsgId int
	Questions               []db.Questions
	StartTime, EndTime      int64
	Log                     []Log
}

type Log struct {
	QuestionId int
	AnswerId   int
}

func (bot *BotApi) ListenForUpdates()  {
	variants = []string{"A) ", "B) ", "C) ", "D) "}

	for update := range bot.Updates {
		bot.Update = update

		if update.Message != nil {
			bot.messageUpdateListener()
		} else {
			if update.CallbackQuery != nil {
				bot.callbackQueryListener()
			}
		}
	}
}

func (bot *BotApi) messageUpdateListener()  {
	if bot.Update.Message.Command() == "" {
		bot.messageListener()
	} else {
		bot.commandListener()
	}
}

func (bot *BotApi) callbackQueryListener()  {
	bot.dynamicCallbackQuery()

	switch bot.Update.CallbackQuery.Data {
	case "schedule":
		schedule := getSchedule()

		msg := newMessage(
			bot.Update.CallbackQuery.Message.Chat.ID,
			getText("schedule") + "\n\n" + schedule.Value,
			"html")

		bot.BotApi.Send(msg)
	case "ask":
		msg := newMessage(
			bot.Update.CallbackQuery.Message.Chat.ID,
			getText("ask"),
			"html")

		bot.BotApi.Send(msg)

		asked = true
	case "faq":
		msg := newMessage(bot.Update.CallbackQuery.Message.Chat.ID, getText("faq"), "html")

		keyboard := tgbotapi.InlineKeyboardMarkup{}
		questions := getFaq()


		for _, item := range questions {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(item.Question, "faq_" + strconv.Itoa(item.Id))
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}

		msg.ReplyMarkup = keyboard
		bot.BotApi.Send(msg)
	case "test":
		msg := newMessage(
			bot.Update.CallbackQuery.Message.Chat.ID,
			getText("test"),
			"html")

		keyboard := tgbotapi.InlineKeyboardMarkup{}

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(getText("startTest") + " " + getEmoji("right-arrow"), "startTest")
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

		msg.ReplyMarkup = keyboard
		bot.BotApi.Send(msg)
	case "startTest":
		bot.BotApi.DeleteMessage(tgbotapi.DeleteMessageConfig{bot.Update.CallbackQuery.Message.Chat.ID, bot.Update.CallbackQuery.Message.MessageID})

		if isExist(bot.Update.CallbackQuery.Message.Chat.UserName) {
			bot.BotApi.DeleteMessage(tgbotapi.DeleteMessageConfig{
				bot.Update.CallbackQuery.Message.Chat.ID,
				Chats[bot.Update.CallbackQuery.Message.Chat.UserName].LastMsgId,
				})
		}

		if checkIfUserExists(bot.Update.CallbackQuery.Message.Chat.UserName) {
			msg := newMessage(
				bot.Update.CallbackQuery.Message.Chat.ID,
				getText("recorded"),
				"html")

			bot.BotApi.Send(msg)
		} else {
			if !isExist(bot.Update.CallbackQuery.Message.Chat.UserName) {
				Chats[bot.Update.CallbackQuery.Message.Chat.UserName] = &Quiz{
					bot.Update.CallbackQuery.Message.Chat.UserName,
					0,
					0,
					0,
					getRandQuestions(),
					time.Now().Unix(),
					0,
					[]Log{},
				}
			}

			bot.newQuestionMessage(bot.Update.CallbackQuery.Message.Chat.ID, bot.Update.CallbackQuery.Message.Chat.UserName)
		}
	}
}

func (bot *BotApi) messageListener() {
	if isExist(bot.Update.Message.Chat.UserName) {
		warningMessage := newMessage(
			bot.Update.Message.Chat.ID,
			getText("continueTest"),
			"html")

		bot.BotApi.Send(warningMessage)
	}

	if asked && bot.Update.Message.Text != "" {
		channelId, err := strconv.ParseInt(config.Toml.Bot.ChannelId, 10, 64);

		if err != nil {
			fmt.Println(err)
		}

		msg := tgbotapi.NewForward(channelId, bot.Update.Message.Chat.ID, bot.Update.Message.MessageID)

		bot.BotApi.Send(msg)

		confirmMsg := newMessage(
			bot.Update.Message.Chat.ID,
			getText("confirmed"),
			"html")

		bot.BotApi.Send(confirmMsg)

		asked = false
	}
}

func (bot *BotApi) commandListener()  {
	switch bot.Update.Message.Command() {
	case "start", "menu":
		msg := newMessage(bot.Update.Message.Chat.ID, "<b>Меню:</b>", "html")

		keyboard := tgbotapi.InlineKeyboardMarkup{}
		menu := getMenu()

		for i, item := range menu {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(item.Name + " " + menuEmojiList()[i], item.Alias)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}

		msg.ReplyMarkup = keyboard
		bot.BotApi.Send(msg)
	}
}

func (bot *BotApi) dynamicCallbackQuery()  {
	if strings.Contains(bot.Update.CallbackQuery.Data, "faq_") {
		bot.faqCallbackQuery()
	} else if strings.Contains(bot.Update.CallbackQuery.Data, "variant_") {
		bot.variantCallbackQuery()
	}
}

func (bot *BotApi) faqCallbackQuery()  {
	s := strings.Split(bot.Update.CallbackQuery.Data, "_")
	id, err := strconv.Atoi(s[1])

	if err != nil {
		fmt.Println(err)
	}

	question := getQuestion(id)

	msg := newMessage(
		bot.Update.CallbackQuery.Message.Chat.ID,
		"<b>" + question.Question + "</b>" + "\n\n" + question.Answer,
		"html")

	bot.BotApi.Send(msg)
}

func (bot *BotApi) variantCallbackQuery() {
	callBackQuery := bot.Update.CallbackQuery
	s := strings.Split(bot.Update.CallbackQuery.Data, "_")
	i, err := strconv.Atoi(s[1])
	id, err := strconv.Atoi(s[2])
	user := callBackQuery.Message.Chat.UserName

	if err != nil {
		fmt.Println(err)
	}

	Chats[user].Log = append(Chats[user].Log, Log{QuestionId: Chats[user].Questions[Chats[user].Index].Id, AnswerId: id})
	Chats[user].Score += Chats[user].Questions[Chats[user].Index].Variants[i].Value

	if Chats[user].Index == 0 {
		Chats[user].Score = Chats[user].Questions[Chats[user].Index].Variants[i].Value
		Chats[user].Index += 1

		bot.newQuestionMessage(callBackQuery.Message.Chat.ID, callBackQuery.Message.Chat.UserName)
	} else if Chats[user].Index == 5 {
		Chats[user].EndTime = time.Now().Unix()

		scoreStr := strconv.Itoa(Chats[user].Score)

		Chats[callBackQuery.Message.Chat.UserName] = Chats[user]

		if err != nil {
			fmt.Println(err)
		}

		if newQuizRecord(Chats[user]) == nil {
			delete(Chats, callBackQuery.Message.Chat.UserName)
		}

		msg := newMessage(
			callBackQuery.Message.Chat.ID,
			getText("score") + scoreStr + "</b>",
			"html")

		bot.BotApi.Send(msg)
	} else {
		Chats[user].Index += 1

		bot.newQuestionMessage(callBackQuery.Message.Chat.ID, callBackQuery.Message.Chat.UserName)
	}

	bot.BotApi.DeleteMessage(tgbotapi.DeleteMessageConfig{callBackQuery.Message.Chat.ID, callBackQuery.Message.MessageID})
}

func (bot *BotApi) newQuestionMessage(chatId int64, userName string) {
	var err error

	indexStr := strconv.Itoa(Chats[userName].Index + 1)
	msg := newMessage(chatId, "<b>" + indexStr + ")</b> " + Chats[userName].Questions[Chats[userName].Index].Text, "html")

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for i, item := range Chats[userName].Questions[Chats[userName].Index].Variants {
		var row []tgbotapi.InlineKeyboardButton

		btn := tgbotapi.NewInlineKeyboardButtonData(variants[i] + item.Text, "variant_" + strconv.Itoa(i) + "_" + strconv.Itoa(item.Id))
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	message, err := bot.BotApi.Send(msg)

	if err != nil {
		fmt.Println(err)
	}

	Chats[userName].LastMsgId = message.MessageID
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

func isExist(user string) bool {
	if _, ok := Chats[user]; ok {
		return true
	}

	return false
}