package main

import (
	"github.com/SkaarjScout/zl-telegram-bot/bothandler"
	"github.com/SkaarjScout/zl-telegram-bot/spreadsheets"
)

type PostgresConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	DbName   string `yaml:"DbName"`
}

type Config struct {
	SpreadsheetsConfig spreadsheets.Config `yaml:"Spreadsheets"`
	TelegramBotConfig  bothandler.Config   `yaml:"TelegramBot"`
	PostgresConfig     PostgresConfig      `yaml:"Postgres"`
}
