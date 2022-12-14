package telegram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Config Config

	mediaGroupChannels map[string]chan orderedMessage
	api                *tgbotapi.BotAPI
	sync.Mutex
}

func (t *Telegram) StartReader() (<-chan Message, error) {
	if err := os.MkdirAll(t.Config.MediaDir, 0750); err != nil {
		return nil, fmt.Errorf("telegram: couldn't create media dir: %s: %w", t.Config.MediaDir, err)
	}

	t.mediaGroupChannels = make(map[string]chan orderedMessage)

	var err error
	t.api, err = tgbotapi.NewBotAPI(t.Config.ApiToken)
	if err != nil {
		return nil, fmt.Errorf("telegram auth failed: %w", err)
	}
	log.Printf("telegram: logged in as %s", t.api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	ch := make(chan Message)
	go func() {
		// If we post 10 images per second, it'd take about 5.8e12 years to max
		// out this counter. This is around 5.8e12 years after the sun explodes
		// and swallows us all.
		var counter uint64 = 0

		for update := range t.api.GetUpdatesChan(updateConfig) {
			update := update
			updateId := counter
			counter++

			if update.ChannelPost == nil {
				continue
			}

			go func() {
				if err := t.handleUpdate(ch, orderedMessage{updateId, update.ChannelPost}); err != nil {
					log.Printf("telegram: handling post: %s", err)
				}
			}()
		}
		close(ch)
	}()

	return ch, nil
}

func (t *Telegram) handleUpdate(ch chan Message, msg orderedMessage) error {
	if msg.MediaGroupID == "" {
		// this message is not part of a group, we do not have to wait
		return t.submitMessages(ch, []orderedMessage{msg})
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

func (t *Telegram) handleMediaGroup(ch chan Message, groupId string) chan orderedMessage {
	groupCh := make(chan orderedMessage)

	go func() {
		messages := messageGroup{}
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
				t.Lock()
				delete(t.mediaGroupChannels, groupId)
				t.Unlock()

				if err := t.submitMessages(ch, messages); err != nil {
					log.Printf("telegram: submitMessages: %s", err)
				}
				return
			}
		}
	}()
	return groupCh
}

type MediaType int

const (
	MediaTypeUnknown MediaType = 0
	MediaTypePhoto   MediaType = iota
	MediaTypeVideo   MediaType = iota
)

func (m MediaType) String() string {
	return []string{"unknown", "photo", "video"}[m]
}

func (t *Telegram) submitMessages(ch chan Message, msgs messageGroup) error {
	if len(msgs) == 0 {
		return nil
	}

	msgs.Sort()

	attachedMedia := []Media{}
	for _, msg := range msgs {
		var mediaType MediaType
		var fileID string
		var uniqueID string
		var caption = msg.Caption
		switch {
		case len(msg.Photo) > 0:
			largest := msg.Photo[0]
			for _, photo := range msg.Photo {
				if photo.Height > largest.Height {
					largest = photo
				}
			}
			fileID = largest.FileID
			uniqueID = largest.FileUniqueID
			mediaType = MediaTypePhoto
		case msg.Video != nil:
			fileID = msg.Video.FileID
			uniqueID = msg.Video.FileUniqueID
			mediaType = MediaTypeVideo
		case msg.Sticker != nil:
			fileID = msg.Sticker.FileID
			uniqueID = msg.Sticker.FileUniqueID
			mediaType = MediaTypePhoto
		default:
			continue
		}

		path, err := t.saveAttachment(fileID, uniqueID)
		if err != nil {
			return fmt.Errorf("telegram: downloading file %s (%s): %w", uniqueID, mediaType, err)
		}

		attachedMedia = append(attachedMedia, Media{
			Path:    path,
			Caption: caption,
		})

	}
	ch <- Message{
		Text:  msgs.Text(),
		Media: attachedMedia,
	}
	return nil
}

func (t *Telegram) saveAttachment(fileID string, uniqueID string) (string, error) {
	url, err := t.api.GetFileDirectURL(fileID)
	if err != nil {
		return "", err
	}
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http error %d", r.StatusCode)
	}
	defer r.Body.Close()

	path := filepath.Join(t.Config.MediaDir, uniqueID)
	f, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	return path, err
}
