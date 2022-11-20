package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mattn/go-mastodon"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type MastodonConfig struct {
	Server       string
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	Token        string `yaml:"token"`
}

func (m *MastodonConfig) CheckAppAuth(ctx context.Context) error {
	if m.Server == "" {
		return fmt.Errorf("no server provided")
	}
	if m.ClientID == "" {
		return fmt.Errorf("no client ID provided")
	}
	if m.ClientSecret == "" {
		return fmt.Errorf("no server ID provided")
	}

	mc := mastodon.NewClient(&mastodon.Config{
		Server:       m.Server,
		ClientID:     m.ClientID,
		ClientSecret: m.ClientSecret,
	})
	return mc.AuthenticateApp(ctx)
}

func (m *MastodonConfig) CheckUserAuth(ctx context.Context) (*mastodon.Account, error) {
	if m.Token == "" {
		return nil, fmt.Errorf("no user token provided")
	}
	mc := mastodon.NewClient(&mastodon.Config{
		Server:       m.Server,
		ClientID:     m.ClientID,
		ClientSecret: m.ClientSecret,
	})
	if err := mc.AuthenticateToken(ctx, m.Token, "urn:ietf:wg:oauth:2.0:oob"); err != nil {
		return nil, err
	}
	return mc.GetAccountCurrentUser(ctx)
}

type TelegramConfig struct {
	ApiToken string   `yaml:"api-token"`
	Channels []string `yaml:"channels"`
}

func (t *TelegramConfig) CheckAuth() (*tgbotapi.User, error) {
	bot, err := tgbotapi.NewBotAPI(t.ApiToken)
	if err != nil {
		return nil, err
	}

	return &bot.Self, nil
}

type Config struct {
	Mastodon MastodonConfig `yaml:"mastodon"`
	Telegram TelegramConfig `yaml:"telegram"`

	path string
}

func LoadConfig(c *Config) error {
	cfgPath, ok := os.LookupEnv("TG2FEDI_CFG")
	if !ok {
		cfgPath = "./tg2fedi.yml"
	}

	cfgF, err := os.Open(cfgPath)
	if err != nil {
		return fmt.Errorf("couldn't open config: %s: %w", cfgPath, err)
	}
	defer cfgF.Close()

	if err = yaml.NewDecoder(cfgF).Decode(&c); err != nil {
		return fmt.Errorf("couldn't decode config: %s: %w", cfgPath, err)
	}

	c.path = cfgPath

	return nil
}

func main() {
	app := &cli.App{
		Name:  "tg2fedi",
		Usage: "",
		Action: func(*cli.Context) error {
			var config Config
			if err := LoadConfig(&config); err != nil {
				log.Printf("to generate a config, run `%s setup`\n", os.Args[0])
				return err
			}
			return nil
		},
		Commands: []*cli.Command{
			SetupCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
