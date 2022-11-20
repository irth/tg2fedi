package config

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/irth/tg2fedi/internal/mastodon"
	mastodonApi "github.com/mattn/go-mastodon"
)

func setupMastodon(ctx context.Context, c *mastodon.Config) error {
	skipAppRegistration := false
	if err := c.CheckAppAuth(ctx); err == nil {
		skipAppRegistration = askBool("Mastodon app client ID/secret already configured. Skip configuring these?", true)
	}

	if !skipAppRegistration {
		if askBool("Do you have an OAuth client ID and secret from your Mastodon instance?", false) {
			if err := setupApp(ctx, c); err != nil {
				return err
			}
		} else {
			if err := createApp(ctx, c); err != nil {
				return err
			}
		}
	}

	skipUser := false
	if u, err := c.CheckUserAuth(ctx); err == nil {
		skipUser = askBool(fmt.Sprintf("User %s (%s) already logged in. Skip?", u.DisplayName, u.Username), true)
	}
	if !skipUser {
		return setupUser(ctx, c)
	}
	return nil
}

func setupApp(ctx context.Context, c *mastodon.Config) error {
	for {
		c.Server = askStr("Your Mastodon instance address (with schema)", c.Server, c.Server)
		c.ClientID = askStr("Client ID", c.ClientID, c.ClientID)
		c.ClientSecret = askStr("Client Secret", c.ClientSecret, c.ClientSecret)
		if err := c.CheckAppAuth(ctx); err != nil {
			fmt.Printf("invalid client id/secret: %s\n", err)
			if askBool("Retry?", true) {
				continue
			}
			return nil
		}
		return nil
	}
}

func createApp(ctx context.Context, c *mastodon.Config) error {
	var server, clientName string
	website := "https://github.com/irth/tg2fedi"

	for {
		server = askStr("Your Mastodon instance address (with schema)", c.Server, c.Server)
		_, err := url.Parse(server)
		if err != nil {
			return fmt.Errorf("invalid server url: %s: %w", server, err)
		}
		clientName = "tg2fedi"
		if hostname, err := os.Hostname(); err == nil {
			clientName = fmt.Sprintf("tg2fedi @ %s", hostname)
		}
		clientName = askStr("A name to identify your app (appears in \"Posted by\")", clientName, clientName)
		website = askStr("Website (linked from \"Posted by\")", website, website)

		fmt.Printf("\nYou're about to create an app with the following configuration:\n\n")
		fmt.Printf("          instance: %s\n", server)
		fmt.Printf("  application name: %s\n", clientName)
		fmt.Printf("            scopes: %s\n", oauthScopes)
		fmt.Printf("           website: %s\n\n", website)

		if !askBool("Continue?", true) {
			continue
		}

		app, err := mastodonApi.RegisterApp(ctx, &mastodonApi.AppConfig{
			Server:     server,
			ClientName: clientName,
			Scopes:     oauthScopes,
			Website:    website,
		})

		if err != nil {
			fmt.Printf("\nCouldn't create the app: %s\n", err)
			if askBool("Do you want to try again?", true) {
				continue
			} else {
				return fmt.Errorf("couldn't create the app: %w", err)
			}
		}

		c.Server = server
		c.ClientID = app.ClientID
		c.ClientSecret = app.ClientSecret

		fmt.Printf("\nApp created on %s:\n\n", c.Server)
		fmt.Printf("      Client ID: %s\n", c.ClientID)
		fmt.Printf("  Client secret: %s\n\n", c.ClientSecret)
		fmt.Println("Keep these secret. You will be given the option to save them to a config file later.")
		return nil
	}
}

func setupUser(ctx context.Context, c *mastodon.Config) error {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return fmt.Errorf("invalid server url: %s: %w", c.Server, err)
	}

	serverURL.Path = "/oauth/authorize"
	q := serverURL.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", "urn:ietf:wg:oauth:2.0:oob")
	q.Set("scope", oauthScopes)
	serverURL.RawQuery = q.Encode()

	fmt.Printf("\nNow, open the following URL and log in to the account that you want the app to post as:\n\n")
	fmt.Println(serverURL.String())
	fmt.Printf("\nAfterwards, paste the token. Keep it secret, as it allows to post as your account.\n")
	for {
		c.Token = askStr("Token", "", "")
		user, err := c.CheckUserAuth(ctx)
		if err != nil {
			fmt.Printf("Invalid token: %s\n", err)
			if askBool("Retry?", true) {
				continue
			}
			return fmt.Errorf("invalid token: %w", err)
		}

		fmt.Printf("Logged in as %s (%s)\n", user.DisplayName, user.Username)
		return nil
	}
}
