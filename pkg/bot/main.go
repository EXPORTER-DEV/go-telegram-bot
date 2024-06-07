package bot

import (
	"time"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/client"
)

type Config struct {
	Token   string
	URL     string
	Limit   int
	Timeout time.Duration
	client  client.Client
}

type ConfigBuilder = func(*Config)

func newConfig(token string) *Config {
	timeout := time.Duration(time.Second * 5)
	return &Config{
		Token: token,
		URL:   "https://api.telegram.org/",
		// Need to double timeout, cause the server can answer a little bit longer cause of high load:
		client:  client.New(client.WithTimeout(timeout * 2)),
		Limit:   50,
		Timeout: timeout,
	}
}

func WithURL(url string) ConfigBuilder {
	return func(c *Config) {
		c.URL = url
	}
}

func WithTimeout(timeout time.Duration) ConfigBuilder {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithLimit(limit int) ConfigBuilder {
	return func(c *Config) {
		c.Limit = limit
	}
}

func New(token string, builders ...ConfigBuilder) (api.Requester, error) {
	config := newConfig(token)

	for _, fn := range builders {
		fn(config)
	}

	return api.New(
		config.Token,
		config.Limit,
		config.URL,
		config.Timeout,
		config.client,
	)
}
