package bothandler

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ConversationManager struct {
	conversationMap map[int64]chan tgbotapi.Update
	stopChan        chan int64
	ctx             context.Context
}

func createConversationManager(ctx context.Context) *ConversationManager {
	return &ConversationManager{
		conversationMap: make(map[int64]chan tgbotapi.Update),
		stopChan:        make(chan int64),
		ctx:             ctx,
	}
}

func (conv *ConversationManager) startConversation(
	convFunc func(
		ctx context.Context,
		chatId int64,
		updateChan <-chan tgbotapi.Update,
		stopChan chan<- int64,
	), update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	updateChan, ok := conv.conversationMap[chatId]
	if !ok {
		updateChan = make(chan tgbotapi.Update)
		conv.conversationMap[chatId] = updateChan
		go convFunc(conv.ctx, chatId, updateChan, conv.stopChan)
	} else {
		return errors.New("conversation already exists, probable bug in handling")
	}
	updateChan <- update
	return nil
}

func (conv *ConversationManager) waitForConversations() {
	for len(conv.conversationMap) > 0 {
		chatId := <-conv.stopChan
		delete(conv.conversationMap, chatId)
	}
}
