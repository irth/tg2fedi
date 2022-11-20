package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Config Config
}

type Message struct {
	Text  string
	Media []string
}

func (t *Telegram) StartReader() (<-chan Message, error) {
	tg, err := tgbotapi.NewBotAPI(t.Config.ApiToken)
	if err != nil {
		return nil, fmt.Errorf("telegram auth failed: %w", err)
	}
	log.Printf("telegram: logged in as %s", tg.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	ch := make(chan Message)
	go func() {
		for msg := range tg.GetUpdatesChan(updateConfig) {
			log.Printf("received: %+v", msg)
			// TODO: support media
			if msg.ChannelPost == nil {
				continue
			}

			m := Message{
				Text: msg.ChannelPost.Text,
			}

			ch <- m
		}
		close(ch)
	}()

	return ch, nil
}
