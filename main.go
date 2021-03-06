package main

import (
	"./src/db"
	"./src/config"
	"./src/handler"
	botPkg "./src/bot"
	"fmt"
	"net/http"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"./src/bot/panel"
)

func main() {
    config.ReadConfigs()

    err := db.OpenConnection()

    if err != nil {
		fmt.Println(err)
    }

    botApi, err := tgbotapi.NewBotAPI(config.Toml.Bot.Token)

    if err != nil {
		fmt.Println(err)
    }

    //botApi.Debug = true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = config.Toml.Bot.Timeout

    updates, err := botApi.GetUpdatesChan(u)

    run()

	bot := botPkg.BotApi{botApi, updates, tgbotapi.Update{}, panel.BotPanel{}}
	bot.ListenForUpdates()
}

func run()  {
	http.HandleFunc("/", handler.MainHandler)
	http.HandleFunc("/health", handler.HealthHandler)
	http.HandleFunc("/upload", handler.FileUploadHandler)
	go http.ListenAndServe(config.Toml.Http.Port, nil)
}