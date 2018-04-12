package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"strconv"
	"strings"
	"time"
	"regexp"
	"../config"
	"../db"
	"./panel"
	"./helper"
)

var variants []string
var Chats = map[int]*Quiz{}
var Asks  = map[int]bool{}
var Panel = map[int64]*PanelSession{}

type PanelSession struct {
	UserId int
	Live   bool
}

type BotApi struct {
	BotApi   *tgbotapi.BotAPI
	Updates  tgbotapi.UpdatesChannel
	Update   tgbotapi.Update
	BotPanel panel.BotPanel
}

type Quiz struct {
	UserId                  int
	UserName                string
	ChatId                  int64
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
	variants = []string{"A", "B", "C", "D"}

	for update := range bot.Updates {
		bot.Update = update
		bot.BotPanel = panel.BotPanel{bot.BotApi, bot.Update}

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
		bot.sendSchedule()
	case "ask":
		bot.sendAskSpeakerMsg()
	case "faq":
		bot.sendFaqMsg()
	case "test":
		bot.sendAboutTestMsg()
	case "startTest":
		bot.sendStartTest()
	}
}

func (bot *BotApi) messageListener() {
	if _, ok := Panel[bot.Update.Message.Chat.ID]; ok && Panel[bot.Update.Message.Chat.ID].Live && bot.BotPanel.IsMessaging() {
		bot.BotPanel.ListenPanelMsgs()
	}

	if isExist(bot.Update.Message.From.ID) {
		warningMessage := helper.NewMessage(
			bot.Update.Message.Chat.ID,
			helper.GetText("continueTest"),
			"html")

		bot.BotApi.Send(warningMessage)
	}

	if ifAsked(bot.Update.Message.From.ID) && bot.Update.Message.Text != "" {
		channelId, err := strconv.ParseInt(config.Toml.Bot.ChannelId, 10, 64);

		if err != nil {
			fmt.Println(err)
		}

		msg := tgbotapi.NewForward(channelId, bot.Update.Message.Chat.ID, bot.Update.Message.MessageID)

		bot.BotApi.Send(msg)

		confirmMsg := helper.NewMessage(
			bot.Update.Message.Chat.ID,
			helper.GetText("confirmed"),
			"html")

		bot.BotApi.Send(confirmMsg)

		Asks[bot.Update.Message.From.ID] = false
	}
}

func (bot *BotApi) commandListener()  {
	cmd := bot.Update.Message.Command()

	switch {
	case regexp.MustCompile("^start$").MatchString(cmd), regexp.MustCompile("^menu$").MatchString(cmd):
		bot.botStartMenu()
	case regexp.MustCompile("^panel$").MatchString(cmd):
		bot.panelStart()
	case regexp.MustCompile("panel_").MatchString(cmd):
		if _, ok := Panel[bot.Update.Message.Chat.ID]; ok && Panel[bot.Update.Message.Chat.ID].Live {
			bot.BotPanel.ListenPanelCmds()
		}
	}
}

func (bot *BotApi) botStartMenu() {
	msg := helper.NewMessage(bot.Update.Message.Chat.ID, "<b>Меню:</b>", "html")

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	menu := getMenu()

	for i, item := range menu {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(item.Name + " " + helper.MenuEmojiList()[i], item.Alias)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	bot.BotApi.Send(msg)
}

func (bot *BotApi) panelStart() {
	if _, ok := Panel[bot.Update.Message.Chat.ID]; ok && Panel[bot.Update.Message.Chat.ID].Live {
		msg := helper.NewMessage(
			bot.Update.Message.Chat.ID,
			"<b>Ваша сессия уже активна, для справки наберите команду</b> /panel_help",
			"html")

		bot.BotApi.Send(msg)
	} else {
		if checkIfAdminExists(bot.Update.Message.From.ID) {
			Panel[bot.Update.Message.Chat.ID] = &PanelSession{bot.Update.Message.From.ID, true}

			msg := helper.NewMessage(
				bot.Update.Message.Chat.ID,
				"<b>Вы вошли в панель управления ботом. Для справки наберите</b> /panel_help",
				"html")

			bot.BotApi.Send(msg)
		} else {
			fmt.Println(fmt.Sprintf("User - %d is trying to signin to panel", bot.Update.Message.From.ID))
		}
	}
}

func (bot *BotApi) dynamicCallbackQuery()  {
	if strings.Contains(bot.Update.CallbackQuery.Data, "faq_") {
		bot.faqCallbackQuery()
	} else if strings.Contains(bot.Update.CallbackQuery.Data, "variant_") {
		bot.variantCallbackQuery()
	} else if strings.Contains(bot.Update.CallbackQuery.Data, "category_") {
		bot.categoryCallbackQuery()
	}
}

func (bot *BotApi) faqCallbackQuery()  {
	s := strings.Split(bot.Update.CallbackQuery.Data, "_")
	id, err := strconv.Atoi(s[1])

	if err != nil {
		fmt.Println(err)
	}

	question := getQuestion(id)

	msg := helper.NewMessage(
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
	user := callBackQuery.From.ID

	if err != nil {
		fmt.Println(err)
	}

	if isExist(user) {
		Chats[user].Log = append(Chats[user].Log, Log{QuestionId: Chats[user].Questions[Chats[user].Index].Id, AnswerId: id})
		Chats[user].Score += Chats[user].Questions[Chats[user].Index].Variants[i].Value

		if Chats[user].Index == 0 {
			Chats[user].Score = Chats[user].Questions[Chats[user].Index].Variants[i].Value
			Chats[user].Index += 1

			bot.newQuestionMessage(callBackQuery.Message.Chat.ID, callBackQuery.From.ID)
		} else if Chats[user].Index == 5 {
			Chats[user].EndTime = time.Now().Unix()

			scoreStr := strconv.Itoa(Chats[user].Score)

			Chats[callBackQuery.From.ID] = Chats[user]

			if err != nil {
				fmt.Println(err)
			}

			if newQuizRecord(Chats[user]) == nil {
				delete(Chats, callBackQuery.From.ID)
			}

			msg := helper.NewMessage(
				callBackQuery.Message.Chat.ID,
				helper.GetText("score") + scoreStr + " очков",
				"html")

			bot.BotApi.Send(msg)
		} else {
			Chats[user].Index += 1

			bot.newQuestionMessage(callBackQuery.Message.Chat.ID, callBackQuery.From.ID)
		}
	} else {
		msg := helper.NewMessage(
			callBackQuery.Message.Chat.ID,
			helper.GetText("testError"),
			"html")

		bot.BotApi.Send(msg)
	}

	bot.BotApi.DeleteMessage(tgbotapi.DeleteMessageConfig{callBackQuery.Message.Chat.ID, callBackQuery.Message.MessageID})
}

func (bot *BotApi) newQuestionMessage(chatId int64, userId int) {
	var (
		err error
		row []tgbotapi.InlineKeyboardButton
	)

	indexStr := strconv.Itoa(Chats[userId].Index + 1)
	questionText := "<b>" + indexStr + ")</b> " + Chats[userId].Questions[Chats[userId].Index].Text + "\n\n"

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for i, item := range Chats[userId].Questions[Chats[userId].Index].Variants {
		btn := tgbotapi.NewInlineKeyboardButtonData(variants[i], "variant_" + strconv.Itoa(i) + "_" + strconv.Itoa(item.Id))
		row = append(row, btn)

		questionText += "<b>" + variants[i] + ")</b> " + item.Text + "\n"
	}

	msg := helper.NewMessage(chatId,  questionText, "html")

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard
	message, err := bot.BotApi.Send(msg)

	if err != nil {
		fmt.Println(err)
	}

	Chats[userId].LastMsgId = message.MessageID
}

func (bot *BotApi) sendSchedule()  {
	schedule := getSchedule()

	msg := helper.NewMessage(
		bot.Update.CallbackQuery.Message.Chat.ID,
		helper.GetText("schedule") + "\n\n" + schedule.Value,
		"html")

	bot.BotApi.Send(msg)
}

func (bot *BotApi) sendAskSpeakerMsg() {
	msg := helper.NewMessage(
		bot.Update.CallbackQuery.Message.Chat.ID,
		helper.GetText("ask"),
		"html")

	bot.BotApi.Send(msg)

	Asks[bot.Update.CallbackQuery.From.ID] = true
}

func (bot *BotApi) sendFaqMsg() {
	msg := helper.NewMessage(bot.Update.CallbackQuery.Message.Chat.ID, helper.GetText("faq"), "html")

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
}

func (bot *BotApi) sendAboutTestMsg() {
	var row []tgbotapi.InlineKeyboardButton

	msg := helper.NewMessage(
		bot.Update.CallbackQuery.Message.Chat.ID,
		helper.GetText("test"),
		"html")

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	btn := tgbotapi.NewInlineKeyboardButtonData(helper.GetText("startTest") + " " + helper.GetEmoji("right-arrow"), "startTest")
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard
	bot.BotApi.Send(msg)
}

func (bot *BotApi) sendStartTest() {
	bot.BotApi.DeleteMessage(tgbotapi.DeleteMessageConfig{bot.Update.CallbackQuery.Message.Chat.ID, bot.Update.CallbackQuery.Message.MessageID})

	if isExist(bot.Update.CallbackQuery.From.ID) {
		bot.BotApi.DeleteMessage(tgbotapi.DeleteMessageConfig{
			bot.Update.CallbackQuery.Message.Chat.ID,
			Chats[bot.Update.CallbackQuery.From.ID].LastMsgId,
		})
	}

	if checkIfUserExists(bot.Update.CallbackQuery.From.ID) {
		msg := helper.NewMessage(
			bot.Update.CallbackQuery.Message.Chat.ID,
			helper.GetText("recorded"),
			"html")

		bot.BotApi.Send(msg)
	} else {
		var row []tgbotapi.InlineKeyboardButton

		categories   := []string{1: "PHP", 2: "JavaScript"}
		text := "<b>Выберите категорию по которой хотите пройти викторину:</b>\n"

		keyboard := tgbotapi.InlineKeyboardMarkup{}

		for i := range categories {
			btn := tgbotapi.NewInlineKeyboardButtonData(categories[i], "category_" + strconv.Itoa(i))
			row = append(row, btn)
		}

		msg := helper.NewMessage(bot.Update.CallbackQuery.Message.Chat.ID,  text, "html")

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

		msg.ReplyMarkup = keyboard
		bot.BotApi.Send(msg)
	}
}

func (bot *BotApi) categoryCallbackQuery() {
	bot.BotApi.DeleteMessage(tgbotapi.DeleteMessageConfig{bot.Update.CallbackQuery.Message.Chat.ID, bot.Update.CallbackQuery.Message.MessageID})

	s := strings.Split(bot.Update.CallbackQuery.Data, "_")
	id, err := strconv.Atoi(s[1])

	if err != nil {
		fmt.Println(err)
	}

	if !isExist(bot.Update.CallbackQuery.From.ID) {
		Chats[bot.Update.CallbackQuery.From.ID] = &Quiz{
			bot.Update.CallbackQuery.From.ID,
			bot.Update.CallbackQuery.From.UserName,
			bot.Update.CallbackQuery.Message.Chat.ID,
			0,
			0,
			0,
			getRandQuestions(id),
			time.Now().Unix(),
			0,
			[]Log{},
		}
	}

	fmt.Println()
	fmt.Println(Chats)
	fmt.Println()

	bot.newQuestionMessage(bot.Update.CallbackQuery.Message.Chat.ID, bot.Update.CallbackQuery.From.ID)
}

func isExist(userId int) bool {
	if _, ok := Chats[userId]; ok {
		return true
	}

	return false
}

func ifAsked(userId int) bool {
	if ask, ok := Asks[userId]; ok {
		return ask
	}

	return false
}