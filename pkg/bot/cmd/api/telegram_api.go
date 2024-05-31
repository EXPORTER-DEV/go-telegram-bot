package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/requests"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/client"
)

type TelegramAPIInterface interface {
	Poll(ctx context.Context)
}

type TelegramAPI struct {
	token      string
	offset     int
	Limit      int
	URL        *url.URL
	Timeout    time.Duration
	client     client.Client
	maxRetries int
}

func (api *TelegramAPI) request(ctx context.Context, method string, body io.Reader) (*http.Response, error) {
	// Copy original URL to patch path for current request below:
	var url url.URL
	url = *api.URL

	// Apply path for current API method:
	url.Path = "bot" + api.token + "/" + method

	// Create new request with context:
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), body)

	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	response, err := api.client.Do(request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *TelegramAPI) getUpdates(ctx context.Context, retry int) error {
	req := requests.NewGetUpdatesRequest(
		api.offset,
		api.Limit,
		api.Timeout.Seconds(),
		[]string{"message"},
	)

	serialized, err := req.Serialize()

	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(serialized)

	res, err := api.request(ctx, requests.GetUpdatesMethod, buf)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Got panic while getUpdates: %v\n", r)
		}

		if res != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		if os.IsTimeout(err) && api.maxRetries > retry {
			retry += 1
			fmt.Printf("Timeouted, going for retry: %v/%v\n", retry, api.maxRetries)
			return api.getUpdates(ctx, retry)
		}

		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Got unexpected http status: %v\n", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	bodyReader, err := res.Request.GetBody()

	requestBody, err := io.ReadAll(bodyReader)

	fmt.Printf("Request body: %s\n", requestBody)

	var response = new(responses.GetUpdatesResponse)

	err = decoder.Decode(response)

	if err != nil {
		return err
	}

	if !response.Ok {
		fmt.Printf("Failed response: %+v\n", response)
		return nil
	}

	fmt.Printf("Response: %+v\n", response)

	if len(response.Result) > 0 {
		// Cause result is in chronic order, so last has the highest updateId:
		api.offset = response.Result[len(response.Result)-1].UpdateId + 1
		fmt.Printf("Updated offset: %+v\n", api.offset)
	}

	return nil
}

func (api *TelegramAPI) Poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Canceled polling cause context done\n")
			return
		default:
			err := api.getUpdates(ctx, 0)
			if err != nil {
				fmt.Printf("Got error: %v\n", err)
			}
		}
	}
}

func New(token string, limit int, rawURL string, timeout time.Duration, client client.Client) (*TelegramAPI, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &TelegramAPI{
		token:      token,
		offset:     0,
		Limit:      limit,
		URL:        url,
		Timeout:    timeout,
		client:     client,
		maxRetries: 5,
	}, nil
}
