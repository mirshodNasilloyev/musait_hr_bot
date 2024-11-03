package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mirshodNasilloyev/musait_hr_bot/pkg/config"
	"github.com/mirshodNasilloyev/musait_hr_bot/pkg/repository"
	"github.com/mirshodNasilloyev/musait_hr_bot/pkg/repository/boltdb"
	"github.com/mirshodNasilloyev/musait_hr_bot/pkg/telegram"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Println("Invalid token", err)
	}
	bot.Debug = true

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	userRepository := boltdb.NewUserRepository(db)
	telegramBot := telegram.NewBot(bot, userRepository, cfg.Messages)
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(repository.ApiKey))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.SpreadSheetId))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		return nil
	})
	return db, nil
}
