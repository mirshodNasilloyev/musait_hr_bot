package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"musaitHrMgBotGo/pkg/config"
	"musaitHrMgBotGo/pkg/repository"
)

var userState = make(map[int64]string)

type Bot struct {
	bot            *tgbotapi.BotAPI
	userRepository repository.UserRepository
	messages       config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, userData repository.UserRepository, messages config.Messages) *Bot {
	return &Bot{bot, userData, messages}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates, err := b.initUpdateChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)
	return nil
}
func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			userID := update.Message.Chat.ID

			//Check user status
			if userState[userID] == "awaiting_api_key" || userState[userID] == "awaiting_spreadsheet_id" {
				fmt.Println(userState[userID])
				b.handleAuthProcess(update)
			} else if update.Message.IsCommand() {
				if err := b.handleCommands(update.Message); err != nil {
					log.Println(err)
				}
			} else {
				if err := b.handleMessage(update.Message); err != nil {
					log.Println(err)
				}
			}
		} else if update.CallbackQuery != nil {
			if err := b.handleCallbackQuery(update); err != nil {
				log.Println(err)
			}
		} else {
			fmt.Errorf("error message type %v", update.Message)
		}
	}
}
func (b *Bot) initUpdateChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u), nil
}
