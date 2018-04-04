package db

import (
	"../config"
	_ "github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type Menu struct {
	Id    int    `db:"id"`
	Alias string `db:"alias"`
	Name  string `db:"name"`
}

type Faq struct {
	Id       int    `db:"id"`
	Question string `db:"question"`
	Answer   string `db:"answer"`
}

type Settings struct {
	Id    int    `db:"id"`
	Alias string `db:"alias"`
	Value string `db:"value"`
}

type Questions struct {
	Id         int    `db:"id"`
	Complexity int    `db:"complexity"`
	Text       string `db:"text"`
	Variants   []Variants
}

type Variants struct {
	Id         int    `db:"id"`
	QuestionId int    `db:"question_id"`
	Text       string `db:"text"`
	Value      int    `db:"value"`
}

type Quiz struct {
	Id        int    `db:"id"`
	User      string `db:"user"`
	Score     int    `db:"score"`
	Log       string `db:"log"`
	StartTime string `db:"start_time"`
	EndTime   string `db:"end_time"`
}

var db *sqlx.DB

func OpenConnection() error {
	var err error

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Toml.Mysql.User,
		config.Toml.Mysql.Password,
		config.Toml.Mysql.Host,
		config.Toml.Mysql.Port,
		config.Toml.Mysql.Database,
	)

	db, err = sqlx.Connect("mysql", connectionString)

	if err == nil {
		err = PingDb()
	}

	return err
}

func CloseConnection() error {
	return db.Close()
}

func PingDb() error {
	err := db.Ping()

	return err
}

func GetMenu() ([]Menu, error) {
	result := []Menu{}

	err := db.Select(&result, `SELECT * FROM menu`)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetFaq() ([]Faq, error) {
	result := []Faq{}

	err := db.Select(&result, "SELECT * FROM faq")

	if err != nil {
		return  result, err
	}

	return result, err
}

func GetQuestion(id int) (Faq, error) {
	result := Faq{}

	err := db.Get(&result, "SELECT * FROM faq WHERE id = ?", id)

	if err != nil {
		return result, err
	}

	return result, err
}

func GetSchedule () (Settings, error) {
	result := Settings{}

	err := db.Get(&result, `SELECT * FROM settings WHERE alias = ?`, "schedule")

	if err != nil {
		return result, err
	}

	return result, err
}

func GetRandomQuestionsByComplexity (limit int, complexity int) ([]Questions, error) {
	result := []Questions{}

	err := db.Select(
		&result,
		"SELECT * FROM questions WHERE complexity = ? ORDER BY RAND() LIMIT 0,?",
		complexity,
		limit)

	if err != nil {
		return result, err
	}

	for i := 0; i < len(result); i++ {
		result[i].Variants, err = GetVariants(result[i].Id)
	}

	return result, err
}

func GetVariants(id int) ([]Variants, error) {
	result := []Variants{}

	err := db.Select(&result, "SELECT * FROM variants WHERE question_id = ? ORDER BY RAND()", id)

	if err != nil {
		return result, err
	}

	return result, err
}