package bot

import (
	"../db"
	"fmt"
)

func getMenu() []db.Categories {
	categories, err := db.GetMenu()

	if err != nil {
		fmt.Println(err)
	}

	return categories
}

func getFaq() []db.Faq  {
	questions, err := db.GetFaq()

	if err != nil {
		fmt.Println(err)
	}

	return questions
}

func getQuestion(id int) db.Faq {
	question, err := db.GetQuestion(id)

	if err != nil {
		fmt.Println(err)
	}

	return  question
}

func getSchedule() db.Settings {
	schedule, err := db.GetSchedule()

	if err != nil {
		fmt.Println(err)
	}

	return schedule
}