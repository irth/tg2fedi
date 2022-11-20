package mastodon

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

type Media struct {
	Path    string
	AltText string
}

type Toot struct {
	Status string
	Media  []Media
}

type Mastodon struct {
	Config Config
}

func (m *Mastodon) StartPoster(ctx context.Context) (chan<- Toot, error) {
	ch := make(chan Toot, 8)
	client := mastodon.NewClient(&mastodon.Config{
		Server:       m.Config.Server,
		ClientID:     m.Config.ClientID,
		ClientSecret: m.Config.ClientSecret,
	})
	err := client.AuthenticateToken(ctx, m.Config.Token, "urn:ietf:wg:oauth:2.0:oob")
	if err != nil {
		return nil, fmt.Errorf("mastodon auth failed: %w", err)
	}
	user, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("mastodon auth failed: %w", err)
	}
	log.Printf("mastodon: logged in as %s (%s), ready to toot", user.DisplayName, user.Username)
	go func() {
		for msg := range ch {
			// TODO: upload media
			// TODO: set toot language
			toot := mastodon.Toot{
				Status: msg.Status,
			}
			status, err := client.PostStatus(ctx, &toot)
			if err != nil {
				log.Printf("mastodon: failed to toot %s: %s", msg, err)
				continue
			}
			log.Printf("mastodon: tooted: %s", status.URL)
		}
	}()
	return ch, nil
}
