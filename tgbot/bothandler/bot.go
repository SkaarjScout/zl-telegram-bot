package bothandler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Find      string = "find"
	Add       string = "add"
	Smalltalk string = "smalltalk"
)

type Bot struct {
	botApi       *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	db           *sql.DB
	ctx          context.Context
	conv         *ConversationManager
	botConfig    Config
}

func New(
	ctx context.Context,
	config Config,
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
		db,
		ctx,
		createConversationManager(ctx),
		config,
	}
}

func formatRowData(row []interface{}) string {
	if row == nil {
		return "Not found"
	}
	return fmt.Sprintf("Ник: %v\nИмя: %v\nБио: %v", row[0], row[1], row[2])
}

func (bot *Bot) StartServe(wg *sync.WaitGroup) {
	updates, err := bot.botApi.GetUpdatesChan(bot.updateConfig)
	if err != nil {
		log.Panicf("Error on creating update channel: %v", err)
	}

	wg.Add(bot.botConfig.WorkerCount)
	for i := 0; i < bot.botConfig.WorkerCount; i++ {
		go bot.processUpdates(bot.ctx, wg, updates)
	}
}

func (bot *Bot) processUpdates(ctx context.Context, wg *sync.WaitGroup, updates <-chan tgbotapi.Update) {
	defer wg.Done()
	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			updateChan, startedConversation := bot.conv.conversationMap[update.Message.Chat.ID]
			switch {
			case startedConversation:
				updateChan <- update
			case !update.Message.IsCommand():
				break
			case update.Message.Command() == Add:
				if err := bot.serveAddUser(ctx, update); err != nil {
					bot.onError(update, err)
				}
			case update.Message.Command() == Smalltalk:
				if err := bot.conv.startConversation(bot.goServeSmalltalk, update); err != nil {
					bot.onError(update, err)
				}
			}
		case <-ctx.Done():
			if bot.botConfig.DebugEnabled {
				log.Print("Goroutine shutting down")
			}
			// graceful shutdown
			bot.conv.waitForConversations()
			return
		case stopChatId := <-bot.conv.stopChan:
			delete(bot.conv.conversationMap, stopChatId)
		}
	}
}

func (bot *Bot) onError(update tgbotapi.Update, err error) {
	bot.botApi.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong"))
	log.Print(err)
}

func (bot *Bot) serveAddUser(ctx context.Context, update tgbotapi.Update) error {
	if err := bot.addUser(ctx, update.Message.From.ID, update.Message.CommandArguments()); err != nil {
		return fmt.Errorf("error on user add: %w", err)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User added")
	if _, err := bot.botApi.Send(msg); err != nil {
		return fmt.Errorf("error on message send: %w", err)
	}
	return nil
}

func (bot *Bot) goServeSmalltalk(
	ctx context.Context,
	chatId int64,
	updateChan <-chan tgbotapi.Update,
	stopChan chan<- int64,
) {
	defer func() {
		if bot.botConfig.DebugEnabled {
			log.Print("Ending smalltalk")
		}
		stopChan <- chatId
	}()
	select {
	case <-ctx.Done():
		return
	case update := <-updateChan:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Smalltalk start")
		if _, err := bot.botApi.Send(msg); err != nil {
			log.Printf("error on message send: %v", err)
			return
		}
	}
	select {
	case <-ctx.Done():
		return
	case update := <-updateChan:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Smalltalk cont")
		if _, err := bot.botApi.Send(msg); err != nil {
			log.Printf("error on message send: %v", err)
			return
		}
	}
	select {
	case <-ctx.Done():
		return
	case update := <-updateChan:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Smalltalk end")
		if _, err := bot.botApi.Send(msg); err != nil {
			log.Printf("error on message send: %v", err)
			return
		}
	}
}
