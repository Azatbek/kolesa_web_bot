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
	index     int
	asked     bool
	started   bool
	questions []db.Questions
	quiz      Quiz
	logs      []Log
)

type Quiz struct {
	User      string
	Score     int
	StartTime string
	EndTime   string
	Log       []Log
}

type Log struct {
	QuestionId int
	AnswerId   int
}

func ListenForUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel)  {
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
	case "test":
		msg := newMessage(
			update.CallbackQuery.Message.Chat.ID,
			"<b>Викторина включает в себя 6 вопросов по 4 варианта. Готовы начинать?</b>",
			"html")

		keyboard := tgbotapi.InlineKeyboardMarkup{}

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData("Начать", "startTest")
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	case "startTest":
		index = 0
		questions = getRandQuestions()

		newQuestionMessage(update.CallbackQuery.Message.Chat.ID, bot)
	}
}

func messageListener(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
	if asked && update.Message.Text != "" {
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
	}
}

func commandListener(update tgbotapi.Update, bot tgbotapi.BotAPI)  {
	switch update.Message.Command() {
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
	s := strings.Split(update.CallbackQuery.Data, "_")
	i, err := strconv.Atoi(s[1])
	id, err := strconv.Atoi(s[2])

	if err != nil {
		fmt.Println(err)
	}

	logs = append(logs, Log{QuestionId: questions[index].Id, AnswerId: id})

	if index == 0 {
		fmt.Println(questions)

		quiz = Quiz{
			User: update.CallbackQuery.Message.From.UserName,
			Score: questions[index].Variants[i].Value,
			StartTime: time.Now().Format("Y-m-d h:i:s"),
		}

		index++
		newQuestionMessage(update.CallbackQuery.Message.Chat.ID, bot)
	} else if index == 5 {
		quiz.Log = logs
		quiz.Score += questions[index].Variants[i].Value

		scoreStr := strconv.Itoa(quiz.Score)

		if err != nil {
			fmt.Println(err)
		}

		msg := newMessage(
			update.CallbackQuery.Message.Chat.ID,
			"<b>Ваш результат: " + scoreStr + "</b>",
			"html")

		bot.Send(msg)
	} else {
		quiz.Score += questions[index].Variants[i].Value

		index++
		newQuestionMessage(update.CallbackQuery.Message.Chat.ID, bot)
	}
}

func newQuestionMessage(chatId int64, bot tgbotapi.BotAPI) {
	msg := newMessage(chatId, "<b>" + questions[index].Text + "</b>", "html")

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for i, item := range questions[index].Variants {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(item.Text, "variant_" + strconv.Itoa(i) + "_" + strconv.Itoa(item.Id))
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	bot.Send(msg)
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