package config

import (
	"context"
	"errors"
	"fmt"
	"os"

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

	for {
		mediaDir := askStr("Where do you want to store media downloaded from Telegram (temporarily)?", c.MediaDir, "./media")
		stat, err := os.Stat(mediaDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				c.MediaDir = mediaDir
				break
			}
			fmt.Printf("Couldn't access the provided directory: %s: %s\n", mediaDir, err)
			continue
		}
		if !stat.IsDir() {
			fmt.Printf("Provided path exists and is not a directory: %s\n", mediaDir)
			continue
		}

		c.MediaDir = mediaDir
		break
	}

	fmt.Println("Configure channel IDs to repost from by modifying the file.")
	return nil
}
