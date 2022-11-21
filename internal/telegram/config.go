package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Config struct {
	ApiToken string `yaml:"api-token"`
	MediaDir string `yaml:"media-dir"`
}

func (t *Config) CheckAuth() (*tgbotapi.User, error) {
	bot, err := tgbotapi.NewBotAPI(t.ApiToken)
	if err != nil {
		return nil, err
	}

	return &bot.Self, nil
}
