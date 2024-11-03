package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken  string
	TelegramBotURL string `mapstructure:"telegram_bot_url"`
	DBPath         string `mapstructure:"db_path"`

	Messages Messages
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	Unautorized  string `mapstructure:"unautorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SavedSuccessfully string `mapstructure:"saved_successfully"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}
	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	cfg.TelegramToken = viper.GetString("token")
	return nil
}
