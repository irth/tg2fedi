package config

import (
	"context"
	"fmt"

	"github.com/irth/tg2fedi/internal/telegram"
)

func setupTelegram(ctx context.Context, c *telegram.Config) error {
	skip := false
	if u, err := c.CheckAuth(); err == nil {
		skip = askBool(fmt.Sprintf("Already logged in to Telegram as %s. Skip configuring Telegram access?", u.UserName), true)
	}
	if !skip {
		return withRetry(
			"Telegram access failed",
			func() error {
				fmt.Println("Message https://t.me/BotFather to create a new bot.")
				c.ApiToken = askStr("Telegram API token", "", "")
				u, err := c.CheckAuth()
				if err != nil {
					return err
				}
				fmt.Printf("Logged in to Telegram as %s.", u.UserName)
				return nil
			},
		)
	}

	fmt.Println("Configure channel IDs to repost from by modifying the file.")
	return nil
}
