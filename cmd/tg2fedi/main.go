package main

import (
	"log"
	"os"

	"github.com/irth/tg2fedi/internal/config"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tg2fedi",
		Usage: "",
		Action: func(*cli.Context) error {
			var cfg config.Config
			if err := config.LoadConfig(&cfg); err != nil {
				log.Printf("to generate a config, run `%s setup`\n", os.Args[0])
				return err
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
