package bot

var text = map[string]string {
	"schedule": "<b>Расписание:</b>",
	"ask": "<b>У вас есть вопрос? Задайте его спикерам</b>",
	"faq": "<b>Часто задаваемые вопросы:</b>",
	"test": "<b>Викторина включает в себя 6 вопросов по 4 варианта. Готовы начинать?</b>",
	"startTest": "Начать",
	"recorded": "<b>Ваш предедущий результат уже записан</b>",
	"continueTest": "<b>Продолжайте викторину!</b>",
	"confirmed": "<b>Ваш вопрос принят!</b>",
	"score": "<b>Ваш результат:</b>",
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