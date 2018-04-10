package bot

import (
	"../db"
	"fmt"
	"encoding/json"
)

func getMenu() []db.Menu {
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

func getRandQuestions() ([]db.Questions) {
	easy, err := db.GetRandomQuestionsByComplexity(3, 0)
	medium, err := db.GetRandomQuestionsByComplexity(2, 1)
	hard, err := db.GetRandomQuestionsByComplexity(1, 2)

	if err != nil {
		fmt.Println(err)
	}

	return append(append(easy, medium...), hard...)
}

func newQuizRecord(quiz *Quiz) error {
	logs, err := json.Marshal(quiz.Log)

	if err != nil {
		fmt.Println(logs)
	}

	recordErr := db.NewQuizRecord(db.Quiz{
		User: quiz.User,
		Score: quiz.Score,
		Log: string(logs),
		StartTime: quiz.StartTime,
		EndTime: quiz.EndTime,
	})

	return recordErr
}

func checkIfUserExists(user string) bool {
	_, err := db.GetUserFromQuiz(user)

	if err != nil {
		return false
	}

	return true
}