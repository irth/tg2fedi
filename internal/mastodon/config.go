package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
)

type Config struct {
	Server       string
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	Token        string `yaml:"token"`
}

func (m *Config) CheckAppAuth(ctx context.Context) error {
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

func (m *Config) CheckUserAuth(ctx context.Context) (*mastodon.Account, error) {
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
