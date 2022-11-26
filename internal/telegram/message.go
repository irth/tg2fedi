package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
)

type Message struct {
	Text  string
	Media []Media
}

type Media struct {
	Path    string
	Caption string
}

type orderedMessage struct {
	no uint64
	*tgbotapi.Message
}

type messageGroup []orderedMessage

func (m messageGroup) Sort() {
	slices.SortFunc(m, func(a orderedMessage, b orderedMessage) bool {
		return a.no < b.no
	})
}

func (msgs messageGroup) Text() string {
	// text vs alt text heuristics because of weird telegram message formats
	if len(msgs) == 0 {
		return ""
	}
	if msgs[0].Text != "" {
		return msgs[0].Text
	}

	// if there was no msg.Text set explicitly, we have the following options:
	//
	// 1. media message with ALL CAPTIONS equal - this should be displayed as
	//    the toot body
	// 2. media message with only first image captioned - telegram then uses the
	//    first caption as .Text
	// 3. media message with multiple different captions - no text, just alt for
	//    each image

	allEqual := true
	othersEmpty := true
	first := msgs[0].Caption
	for _, m := range msgs[1:] {
		if m.Caption != first {
			allEqual = false
		}
		if m.Caption != "" {
			othersEmpty = false
		}
	}

	if allEqual || othersEmpty {
		return first
	}
	return ""
}
