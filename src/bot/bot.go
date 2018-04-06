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
	message   tgbotapi.Message
	index     int
	asked     bool
	started   bool
	lastQId   int
	questions []db.Questions
	quiz      Quiz
	logs      []Log
	variants  []string
)

type Quiz struct {
	User      string
	Score     int
	StartTime int64
	EndTime   int64
	Log       []Log
}

type Log struct {
	QuestionId int
	AnswerId   int
}

func ListenForUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel)  {
	variants = []string{"A) ", "B) ", "C) ", "D) "}

	for update := range updates {
		if update.Message != nil {
			messageUpdateListener(update, *bot)
		} else {
			if update.CallbackQuery != nil {
				callbackQueryListener(update, *bot)
			}
		}
	}
}

func messageUpdateListener(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
	if update.Message.Command() == "" {
		messageListener(update, bot)
	} else {
		commandListener(update, bot)
	}
}

func callbackQueryListener(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
	dynamicCallbackQuery(update, bot)

	switch update.CallbackQuery.Data {
	case "schedule":
		schedule := getSchedule()

		msg := newMessage(
			update.CallbackQuery.Message.Chat.ID,
			getText("schedule") + "\n\n" + schedule.Value,
			"html")

		bot.Send(msg)
	case "ask":
		msg := newMessage(
			update.CallbackQuery.Message.Chat.ID,
			getText("ask"),
			"html")

		bot.Send(msg)

		asked = true
	case "faq":
		msg := newMessage(update.CallbackQuery.Message.Chat.ID, getText("faq"), "html")

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
	case "test":
		msg := newMessage(
			update.CallbackQuery.Message.Chat.ID,
			getText("test"),
			"html")

		keyboard := tgbotapi.InlineKeyboardMarkup{}

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(getText("startTest") + " " + getEmoji("right-arrow"), "startTest")
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	case "startTest":
		bot.DeleteMessage(tgbotapi.DeleteMessageConfig{update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID})

		if started {
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{update.CallbackQuery.Message.Chat.ID, lastQId})
		}

		if checkIfUserExists(update.CallbackQuery.Message.Chat.UserName) {
			msg := newMessage(
				update.CallbackQuery.Message.Chat.ID,
				getText("recorded"),
				"html")

			bot.Send(msg)
		} else {
			index = 0
			questions = getRandQuestions()
			started = true

			newQuestionMessage(update.CallbackQuery.Message.Chat.ID, bot)
		}
	}
}

func messageListener(update tgbotapi.Update, bot tgbotapi.BotAPI) {
	if started {
		warningMessage := newMessage(
			update.Message.Chat.ID,
			getText("continueTest"),
			"html")

		bot.Send(warningMessage)
	}

	if asked && update.Message.Text != "" {
		channelId, err := strconv.ParseInt(config.Toml.Bot.ChannelId, 10, 64);

		if err != nil {
			fmt.Println(err)
		}

		msg := tgbotapi.NewForward(channelId, update.Message.Chat.ID, update.Message.MessageID)

		bot.Send(msg)

		confirmMsg := newMessage(
			update.Message.Chat.ID,
			getText("confirmed"),
			"html")

		bot.Send(confirmMsg)

		asked = false
	}
}

func commandListener(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
	switch update.Message.Command() {
	case "start", "menu":
		msg := newMessage(update.Message.Chat.ID, "<b>Меню:</b>", "html")

		keyboard := tgbotapi.InlineKeyboardMarkup{}
		menu := getMenu()

		for i, item := range menu {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(item.Name + " " + menuEmojiList()[i], item.Alias)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}

		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}

func dynamicCallbackQuery(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
	if strings.Contains(update.CallbackQuery.Data, "faq_") {
		faqCallbackQuery(update, bot)
	} else if strings.Contains(update.CallbackQuery.Data, "variant_") {
		variantCallbackQuery(update, bot)
	}
}

func faqCallbackQuery(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
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

func variantCallbackQuery(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
	callBackQuery := update.CallbackQuery
	s := strings.Split(update.CallbackQuery.Data, "_")
	i, err := strconv.Atoi(s[1])
	id, err := strconv.Atoi(s[2])

	if err != nil {
		fmt.Println(err)
	}

	logs = append(logs, Log{QuestionId: questions[index].Id, AnswerId: id})

	if index == 0 {
		quiz = Quiz{
			User: callBackQuery.Message.Chat.UserName,
			Score: questions[index].Variants[i].Value,
			StartTime: time.Now().Unix(),
		}

		index++
		newQuestionMessage(callBackQuery.Message.Chat.ID, bot)
	} else if index == 5 {
		quiz.Log = logs
		quiz.Score += questions[index].Variants[i].Value
		quiz.EndTime = time.Now().Unix()

		scoreStr := strconv.Itoa(quiz.Score)

		if err != nil {
			fmt.Println(err)
		}

		started = false
		newQuizRecord(quiz)

		msg := newMessage(
			callBackQuery.Message.Chat.ID,
			getText("score") + scoreStr + "</b>",
			"html")

		bot.Send(msg)
	} else {
		quiz.Score += questions[index].Variants[i].Value

		index++
		newQuestionMessage(callBackQuery.Message.Chat.ID, bot)
	}

	bot.DeleteMessage(tgbotapi.DeleteMessageConfig{callBackQuery.Message.Chat.ID, callBackQuery.Message.MessageID})
}

func newQuestionMessage(chatId int64, bot tgbotapi.BotAPI) {
	var err error
	indexStr := strconv.Itoa(index + 1)
	msg := newMessage(chatId, "<b>" + indexStr + ")</b> " + questions[index].Text, "html")

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for i, item := range questions[index].Variants {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(variants[i] + item.Text, "variant_" + strconv.Itoa(i) + "_" + strconv.Itoa(item.Id))
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	message, err = bot.Send(msg)

	if err != nil {
		fmt.Println(err)
	}

	lastQId = message.MessageID
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