package helper

import "github.com/go-telegram-bot-api/telegram-bot-api"

var text = map[string]string {
	"schedule": "<b>Расписание конференции:</b>",
	"ask": "<b>У вас есть вопрос? Задайте его спикерам</b>",
	"faq": "<b>Часто задаваемые вопросы:</b>",
	"test": "Викторина включает в себя 6 вопросов с 4 вариантами ответов.\nВремя на  прохождение: 7 мин.\nУ тебя будет одна попытка.\nГотовы начинать?",
	"startTest": "Начинать",
	"recorded": "<b>Ваш предедущий результат уже записан. Подождите подведения итогов</b> /menu",
	"continueTest": "Это напоминалка о том, что вы чатитесь со мной чем продолжать викторину! /menu",
	"confirmed": "<b>Ваш вопрос принят!</b> /menu",
	"scoreTest": "Спасибо большое за участие в нашей викторине.\nСейчас мы будем подводить итоги и ты получишь сообщение с результатом.",
	"testError": "Упс! Что-то пошло не так, попробуйте заново пройти викторину /menu",
	"category": "Здорово, что ты решил поучаствовать в нашей Викторине, ты уже молодец)\nТебе нужно выбрать, в чем ты круче: PHP или Java Script?\nОтветить на наши вопросы до 17:00 в удобное для тебя время (в обед либо во время кофе брейка)\nНа закрытии  конференции мы подведем итоги)\nУдачи!",
	"notEnoughTime": "Спасибо за участие в нашей викторине, к сожалению тебе не хватило времени.\nУверены, что в следующий раз удача улыбнется тебе)",
	"getPrize": "Ура! Поздравляем! Ты круто ответил на вопросы, на закрытие конференции с нас памятный подарочек)",
	"nexTime": "Спасибо, за участие в нашей викторине. К сожалению вы не набрали достаточное количество баллов, для получения подарка от нас. Уверены, что в следующий раз удача улыбнется тебе)",
}

var emoji = map[string]string {
	"calendar": "\xF0\x9F\x93\x85",
	"100point": "\xF0\x9F\x92\xAF",
	"speech": "\xF0\x9F\x92\xAC",
	"question": "\xE2\x9D\x93",
	"right-arrow": "\xE2\x9E\xA1",
}

func GetText(key string) string {
	return text[key]
}

func GetEmoji(key string) string {
	return emoji[key]
}

func MenuEmojiList() []string {
	return []string{
		GetEmoji("calendar"),
		GetEmoji("100point"),
		GetEmoji("speech"),
		GetEmoji("question"),
	}
}

func NewMessage(chatId int64, text string, parseMode string) tgbotapi.MessageConfig {
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