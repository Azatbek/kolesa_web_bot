package db

import (
	"../config"
	_ "github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type Categories struct {
	Id    int    `db:"id"`
	Alias string `db:"alias"`
	Name  string `db:"name"`
}

type Faq struct {
	Id       int `db:"id"`
	Question string `db:"question"`
	Answer   string `db:"answer"`
}

type Settings struct {
	Id    int    `db:"id"`
	Alias string `db:"alias"`
	Value string `db:"value"`
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

func GetMenu() ([]Categories, error) {
	result := []Categories {}

	err := db.Select(&result, `SELECT * FROM menu`)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetFaq() ([]Faq, error) {
	result := []Faq {}

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

	fmt.Println(err)

	if err != nil {
		return result, err
	}

	return result, err
}