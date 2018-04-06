package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
    Http struct {
        Port string
    }

    Mysql struct {
        Host     string
		Port     int
		User     string
		Password string
		Database string
    }

    Bot struct {
		Token     string
        Timeout   int
        ChannelId string
    }
}

var Toml TomlConfig

func ReadConfigs() {
    if _, err := toml.DecodeFile("src/config/config.toml", &Toml); err != nil {
        fmt.Println("Could not read config file")
	fmt.Println(err)
    }
}