package telegram

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestMessageGroupText(t *testing.T) {
	i := uint64(0)
	m := func(text string, caption string) orderedMessage {
		msg := orderedMessage{no: i, Message: &tgbotapi.Message{
			Text:    text,
			Caption: caption,
		}}
		i += 1
		return msg
	}

	testPatterns := []struct {
		expectedTitle string
		msgs          messageGroup
	}{
		{"", messageGroup{}},
		{"hello", messageGroup{m("hello", "")}},
		{"hello", messageGroup{m("hello", ""), m("", "hello")}},
		{"hello", messageGroup{m("", "hello"), m("", "hello"), m("", "hello")}},
		{"", messageGroup{m("", "hello"), m("", "hello"), m("", "helo")}},
		{"hello", messageGroup{m("hello", "a"), m("", "a"), m("", "a")}},
		{"hello", messageGroup{m("hello", "a"), m("", ""), m("", "")}},
		{"hello", messageGroup{m("", "hello"), m("", ""), m("", "")}},
		{"a", messageGroup{m("a", "hello"), m("", ""), m("", "")}},
	}

	for _, tp := range testPatterns {
		assert.Equal(t, tp.expectedTitle, tp.msgs.Text())
	}
}
