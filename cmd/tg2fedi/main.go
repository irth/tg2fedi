package main

import (
	"log"
	"os"

	"github.com/irth/tg2fedi/internal/config"
	"github.com/irth/tg2fedi/internal/mastodon"
	"github.com/irth/tg2fedi/internal/telegram"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tg2fedi",
		Usage: "",
		Action: func(ctx *cli.Context) error {
			var cfg config.Config
			if err := config.LoadConfig(&cfg); err != nil {
				log.Printf("to generate a config, run `%s setup`\n", os.Args[0])
				return err
			}

			m := mastodon.Mastodon{Config: cfg.Mastodon}
			toots, err := m.StartPoster(ctx.Context)
			if err != nil {
				return err
			}
			defer close(toots)

			t := telegram.Telegram{Config: cfg.Telegram}
			posts, err := t.StartReader()
			if err != nil {
				return err
			}

			for post := range posts {
				toots <- mastodon.Toot{
					Status: post.Text,
				}
			}
			return nil
		},
		Commands: []*cli.Command{
			config.SetupCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
