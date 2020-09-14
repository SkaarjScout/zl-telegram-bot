package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/SkaarjScout/zl-telegram-bot/spreadsheets"
)

var TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
var SPREADSHEET_ID = os.Getenv("SPREADSHEET_ID")
var SHEETS_REFRESH_TOKEN = os.Getenv("SHEETS_REFRESH_TOKEN")

func formatRowData(row []interface{}) string {
	if row == nil {
		return "Not found"
	}
	return fmt.Sprintf("Ник: %v\nИмя: %v\nБио: %v", row[0], row[1], row[2])
}

func main() {
	credentialsJson, err := ioutil.ReadFile("credentials.json")
	sheetsClient := spreadsheets.NewClient(spreadsheets.Config{
		SpreadsheetId:   SPREADSHEET_ID,
		CredentialsJson: string(credentialsJson),
		RefreshToken:    SHEETS_REFRESH_TOKEN,
	})

	bot, err := tgbotapi.NewBotAPI(TELEGRAM_TOKEN)
	if err != nil {
		log.Fatalf("Unable to connect to telegram bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	for update := range updates {
		switch {
		case !update.Message.IsCommand():
			continue
		case update.Message.Command() == "find":
			nickname := update.Message.CommandArguments()
			row, err := sheetsClient.FindRow(nickname)
			if err != nil {
				log.Printf("Error on row find: %v", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatRowData(row))
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error on message send: %v", err)
			}
		}

	}
}
