package db

import (
	"../config"
	_ "github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type Categories struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
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

func GetCategories() ([]Categories, error){
	result := []Categories {}

	err := db.Select(&result, `SELECT * FROM categories`)

	if err != nil {
		return result, err
	}

	return result, nil
}