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

func getRandQuestions(category int) ([]db.Questions) {
	easy, err := db.GetRandomQuestionsByComplexity(3, 0, category)
	medium, err := db.GetRandomQuestionsByComplexity(2, 1, category)
	hard, err := db.GetRandomQuestionsByComplexity(1, 2, category)

	if err != nil {
		fmt.Println(err)
	}

	return append(append(easy, hard...), medium...)
}

func newQuizRecord(quiz *Quiz) error {
	logs, err := json.Marshal(quiz.Log)

	if err != nil {
		fmt.Println(logs)
	}

	recordErr := db.NewQuizRecord(db.Quiz{
		UserId: quiz.UserId,
		UserName: quiz.UserName,
		ChatId: quiz.ChatId,
		Score: quiz.Score,
		Log: string(logs),
		StartTime: quiz.StartTime,
		EndTime: quiz.EndTime,
	})

	return recordErr
}

func checkIfUserExists(userId int) bool {
	_, err := db.GetUserFromQuiz(userId)

	fmt.Println(err)

	if err != nil {
		return false
	}

	return true
}

func checkIfAdminExists(userId int) bool {
	_, err := db.GetAdmin(userId)

	fmt.Println(err)

	if err != nil {
		return false
	}

	return true
}