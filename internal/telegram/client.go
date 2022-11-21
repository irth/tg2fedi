package telegram

import (
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Config Config

	mediaGroupChannels map[string]chan *tgbotapi.Message
	sync.Mutex
}

type Message struct {
	Text  string
	Media []string
}

func (t *Telegram) StartReader() (<-chan Message, error) {
	t.mediaGroupChannels = make(map[string]chan *tgbotapi.Message)

	tg, err := tgbotapi.NewBotAPI(t.Config.ApiToken)
	if err != nil {
		return nil, fmt.Errorf("telegram auth failed: %w", err)
	}
	log.Printf("telegram: logged in as %s", tg.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	ch := make(chan Message)
	go func() {
		for update := range tg.GetUpdatesChan(updateConfig) {
			update := update

			if update.ChannelPost == nil {
				continue
			}

			go func() {
				if err := t.handleUpdate(ch, update.ChannelPost); err != nil {
					log.Printf("telegram: handling post: %s", err)
				}
			}()
		}
		close(ch)
	}()

	return ch, nil
}

func (t *Telegram) handleUpdate(ch chan Message, msg *tgbotapi.Message) error {
	if msg.MediaGroupID == "" {
		// this message is not part of a group, we do not have to wait
		return t.submitMessages(ch, []*tgbotapi.Message{msg})
	}

	t.Lock()
	defer t.Unlock()

	groupCh, ok := t.mediaGroupChannels[msg.MediaGroupID]
	if !ok {
		groupCh = t.handleMediaGroup(ch, msg.MediaGroupID)
		t.mediaGroupChannels[msg.MediaGroupID] = groupCh
	}
	groupCh <- msg

	return nil
}

func (t *Telegram) handleMediaGroup(ch chan Message, groupId string) chan *tgbotapi.Message {
	groupCh := make(chan *tgbotapi.Message)

	go func() {
		messages := []*tgbotapi.Message{}
		timerDuration := 10 * time.Second
		timer := time.NewTimer(timerDuration)

		for {
			select {
			case msg := <-groupCh:
				messages = append(messages, msg)
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(timerDuration)
			case <-timer.C:
				func() {
					t.Lock()
					defer t.Unlock()
					delete(t.mediaGroupChannels, groupId)
				}()
				if err := t.submitMessages(ch, messages); err != nil {
					log.Printf("telegram: submitMessages: %s", err)
				}
				return
			}
		}
	}()
	return groupCh
}

func (t *Telegram) submitMessages(ch chan Message, msgs []*tgbotapi.Message) error {
	log.Printf("received msgs: %d", len(msgs))
	for i, msg := range msgs {
		fmt.Println(i, msg.Text, msg.Photo, msg.Caption)
	}
	return nil
}
