package client

import (
	"net/http"
	"time"
)

//go:generate mockery --name Client
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	Timeout time.Duration
}

type ClientConfigBuilder = func(cc *ClientConfig)

func newClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout: time.Duration(time.Second * 30),
	}
}

func WithTimeout(duration time.Duration) ClientConfigBuilder {
	return func(cc *ClientConfig) {
		cc.Timeout = duration
	}
}

func New(builders ...ClientConfigBuilder) Client {
	config := newClientConfig()

	for _, cc := range builders {
		cc(config)
	}

	return &http.Client{
		Timeout: config.Timeout,
	}
}
