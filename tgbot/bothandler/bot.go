package bothandler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/SkaarjScout/zl-telegram-bot/spreadsheets"
)

const (
	Find string = "find"
	Add  string = "add"
)

type Bot struct {
	botApi             *tgbotapi.BotAPI
	updateConfig       tgbotapi.UpdateConfig
	spreadsheetsClient *spreadsheets.Client
	db                 *sql.DB
	botConfig          Config
}

func New(
	config Config,
	spreadsheetsClient *spreadsheets.Client,
	db *sql.DB) Bot {
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
		db,
		config,
	}
}

func formatRowData(row []interface{}) string {
	if row == nil {
		return "Not found"
	}
	return fmt.Sprintf("Ник: %v\nИмя: %v\nБио: %v", row[0], row[1], row[2])
}

func (bot *Bot) StartServe(ctx context.Context, wg *sync.WaitGroup) {
	updates, err := bot.botApi.GetUpdatesChan(bot.updateConfig)
	if err != nil {
		log.Panicf("Error on creating update channel: %v", err)
	}

	wg.Add(bot.botConfig.WorkerCount)
	for i := 0; i < bot.botConfig.WorkerCount; i++ {
		go bot.processUpdates(ctx, wg, updates)
	}
}

func (bot *Bot) processUpdates(ctx context.Context, wg *sync.WaitGroup, updates <-chan tgbotapi.Update) {
	defer wg.Done()
	for {
		select {
		case update := <-updates:
			switch {
			case update.Message == nil || !update.Message.IsCommand():
				break
			case update.Message.Command() == Find:
				if err := bot.serveFind(update); err != nil {
					log.Print(err)
				}
			case update.Message.Command() == Add:
				if err := bot.serveAddUser(ctx, update); err != nil {
					log.Print(err)
				}
			}
		case <-ctx.Done():
			if bot.botConfig.DebugEnabled {
				log.Print("Goroutine shutting down")
			}
			// graceful shutdown
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

func (bot *Bot) serveAddUser(ctx context.Context, update tgbotapi.Update) error {
	if err := bot.addUser(ctx, update.Message.From.ID); err != nil {
		return fmt.Errorf("error on user add: %w", err)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User added")
	if _, err := bot.botApi.Send(msg); err != nil {
		return fmt.Errorf("error on message send: %w", err)
	}
	return nil
}
