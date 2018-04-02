export GOPATH=$(CURDIR)/.go

APP_NAME = kolesa-conf-bot
VERSION = `$(CURDIR)/out/$(APP_NAME) -v | cut -d ' ' -f 3`

$(CURDIR)/out/$(APP_NAME): $(CURDIR)/main.go
	go build -o $(CURDIR)/out/$(APP_NAME) $(CURDIR)/main.go

dep-install:
	go get github.com/BurntSushi/toml
	go get github.com/codegangsta/cli
	go get github.com/endeveit/go-snippets/cli
	go get github.com/go-telegram-bot-api/telegram-bot-api
	go get github.com/jmoiron/sqlx

run:
	go run $(CURDIR)/src/main.go