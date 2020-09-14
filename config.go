package main

import (
	"github.com/SkaarjScout/zl-telegram-bot/spreadsheets"
	"github.com/SkaarjScout/zl-telegram-bot/tgbot"
)

type Config struct {
	SpreadsheetsConfig spreadsheets.Config `yaml:"Spreadsheets"`
	TelegramBotConfig  tgbot.Config        `yaml:"TelegramBot"`
}
