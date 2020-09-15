package main

import (
	"github.com/SkaarjScout/zl-telegram-bot/bothandler"
	"github.com/SkaarjScout/zl-telegram-bot/spreadsheets"
)

type Config struct {
	SpreadsheetsConfig spreadsheets.Config `yaml:"Spreadsheets"`
	TelegramBotConfig  bothandler.Config   `yaml:"TelegramBot"`
}
