package config

import (
	"github.com/irth/tg2fedi/internal/mastodon"
	"github.com/irth/tg2fedi/internal/telegram"
)

type Config struct {
	Mastodon mastodon.Config `yaml:"mastodon"`
	Telegram telegram.Config `yaml:"telegram"`
}
