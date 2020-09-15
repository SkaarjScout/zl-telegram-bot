package bothandler

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/SkaarjScout/zl-telegram-bot/spreadsheets"
)

const (
	Find string = "find"
)

type Bot struct {
	botApi             *tgbotapi.BotAPI
	updateConfig       tgbotapi.UpdateConfig
	spreadsheetsClient *spreadsheets.Client
	botConfig          Config
}

func New(config Config, spreadsheetsClient *spreadsheets.Client) Bot {
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Panicf("Unable to connect to telegram bot: %v", err)
	}

	bot.Debug = config.DebugEnabled
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = config.UpdateTimeout
	return Bot{
		bot,
		updateConfig,
		spreadsheetsClient,
		config,
	}
}

func formatRowData(row []interface{}) string {
	if row == nil {
		return "Not found"
	}
	return fmt.Sprintf("Ник: %v\nИмя: %v\nБио: %v", row[0], row[1], row[2])
}

func (bot *Bot) Serve(stop chan bool) {
	updates, err := bot.botApi.GetUpdatesChan(bot.updateConfig)
	if err != nil {
		log.Panicf("Error on creating update channel: %v", err)
	}

	for {
		select {
		case update := <-updates:
			switch {
			case update.Message == nil || !update.Message.IsCommand():
				break
			case update.Message.Command() == Find:
				log.Print("Serving find")
				if err = bot.serveFind(update); err != nil {
					log.Print(err)
				}
			}
		case <-stop:
			log.Print("Got a stop signal")
			stop <- true
			return
		}
	}
}

func (bot *Bot) serveFind(update tgbotapi.Update) error {
	row, err := bot.spreadsheetsClient.FindRow(update.Message.CommandArguments())
	if err != nil {
		return fmt.Errorf("error on row find: %w", err)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatRowData(row))
	if _, err := bot.botApi.Send(msg); err != nil {
		return fmt.Errorf("error on message send: %w", err)
	}
	return nil
}
