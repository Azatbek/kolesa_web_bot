package bot

var text = map[string]string {
	"schedule": "<b>Расписание:</b>",
	"ask": "<b>У вас есть вопрос? Задайте его спикерам</b>",
	"faq": "<b>Часто задаваемые вопросы:</b>",
	"test": "<b>Викторина включает в себя 6 вопросов по 4 варианта. Готовы начинать?</b>",
	"startTest": "Начать",
	"recorded": "<b>Ваш предедущий результат уже записан</b> /menu",
	"continueTest": "Это напоминалка о том, что вы чатитесь со мной чем продолжать викторину! /menu",
	"confirmed": "<b>Ваш вопрос принят!</b> /menu",
	"score": "<b>Ваш результат:</b> ",
	"testError": "Упс! Что-то пошло не так, попробуйте заново пройти викторину /menu",
}

var emoji = map[string]string {
	"calendar": "\xF0\x9F\x93\x85",
	"100point": "\xF0\x9F\x92\xAF",
	"speech": "\xF0\x9F\x92\xAC",
	"question": "\xE2\x9D\x93",
	"right-arrow": "\xE2\x9E\xA1",
}

func getText(key string) string {
	return text[key]
}

func getEmoji(key string) string {
	return emoji[key]
}

func menuEmojiList() []string {
	return []string{
		getEmoji("calendar"),
		getEmoji("100point"),
		getEmoji("speech"),
		getEmoji("question"),
	}
}