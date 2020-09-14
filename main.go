package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SkaarjScout/zl-telegram-bot/spreadsheets"
	"github.com/SkaarjScout/zl-telegram-bot/tgbot"
)

var TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
var SPREADSHEET_ID = os.Getenv("SPREADSHEET_ID")
var SHEETS_REFRESH_TOKEN = os.Getenv("SHEETS_REFRESH_TOKEN")

func main() {
	credentialsJson, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Panicf("Error on credentials file read: %v", err)
	}
	config := Config{
		SpreadsheetsConfig: spreadsheets.Config{
			SpreadsheetId:   SPREADSHEET_ID,
			CredentialsJson: string(credentialsJson),
			RefreshToken:    SHEETS_REFRESH_TOKEN,
		},
		TelegramBotConfig: tgbot.Config{
			TelegramBotToken: TELEGRAM_TOKEN,
			DebugEnabled:     true,
			UpdateTimeout:    60,
		},
	}
	spreadsheetsClient := spreadsheets.NewClient(config.SpreadsheetsConfig)
	bot := tgbot.New(config.TelegramBotConfig, &spreadsheetsClient)

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	botStop := make(chan bool, 1)
	go bot.Serve(botStop)

	select {
	case <-interrupt:
		botStop <- true
	}
}
