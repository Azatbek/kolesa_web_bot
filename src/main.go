package src

import (
	"./db"
	"./config"
	"./handler"
	"./bot"
	"fmt"
	"net/http"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
    config.ReadConfigs()

    err := db.OpenConnection()

    if err != nil {
	fmt.Println(err)
    }
    botApi, err := tgbotapi.NewBotAPI(config.Toml.Bot.Token)

    if err != nil {
	fmt.Println("Bot connection failed")
    }

    botApi.Debug = true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = config.Toml.Bot.Timeout

    updates, err := botApi.GetUpdatesChan(u)

    run()

    fmt.Println("here")

    bot.ListenForUpdates(botApi, updates)
}

func run()  {
	http.HandleFunc("/", handler.MainHandler)
	http.HandleFunc("/health", handler.HealthHandler)
	go http.ListenAndServe(config.Toml.Http.Port, nil)
}