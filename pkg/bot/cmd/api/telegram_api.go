package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/domain"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/requests"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/client"
)

var ErrInvalidResponse = errors.New("INVALID_TELEGRAM_API_RESPONSE")
var ErrInvalidArgument = errors.New("INVALID_ARGUMENT")

type TelegramAPIInterface interface {
	Poll(ctx context.Context) chan *responses.Update
	SendMessage(ctx context.Context, message *domain.MessageBuilder) error
	Reply(ctx context.Context, update *responses.Update, text string) error
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

type method string

var getUpdatesMethod method = "getUpdates"
var sendMessageMethod method = "sendMessage"

func (api *TelegramAPI) request(ctx context.Context, method method, body io.Reader) (*http.Response, error) {
	// Copy original URL to patch path for current request below:
	var url url.URL = *api.URL

	// Apply path for current API method:
	url.Path = "bot" + api.token + "/" + string(method)

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

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("UNEXPECTED HTTP STATUS: %v", response.StatusCode)
	}

	return response, nil
}

func (api *TelegramAPI) getUpdates(ctx context.Context, retry int) ([]*responses.Update, error) {
	req := requests.NewGetUpdatesRequest(
		api.offset,
		api.Limit,
		api.Timeout.Seconds(),
		[]string{"message"},
	)

	serialized, err := req.Serialize()

	if err != nil {
		return nil, fmt.Errorf("SERIALIZE_FAILED: %w", err)
	}

	buf := bytes.NewBuffer(serialized)

	res, err := api.request(ctx, getUpdatesMethod, buf)

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

		return nil, err
	}

	decoder := json.NewDecoder(res.Body)

	// DEBUG CODE BLOCK:
	bodyReader, _ := res.Request.GetBody()

	requestBody, _ := io.ReadAll(bodyReader)

	fmt.Printf("Request body: %s\n", requestBody)
	// DEBUG CODE BLOCK END;

	var response = new(responses.GetUpdatesResponse)

	err = decoder.Decode(response)

	if err != nil {
		return nil, err
	}

	if !response.Ok {
		fmt.Printf("Got incorrect response from Telegram API: %+v", response)
		return nil, fmt.Errorf("GOT_INCORRECT_RESPONSE: %w", ErrInvalidResponse)
	}

	fmt.Printf("Response: %+v\n", response)

	if len(response.Result) > 0 {
		// Cause result is in chronic order, so last has the highest updateId:
		api.offset = response.Result[len(response.Result)-1].UpdateId + 1
		fmt.Printf("Updated offset: %+v\n", api.offset)
	}

	return response.Result, nil
}

func (api *TelegramAPI) Poll(ctx context.Context) chan *responses.Update {
	res := make(chan *responses.Update)

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Canceled polling cause context done\n")
				// Probably extra operation, but let it exist here:
				close(res)
				return
			default:
				updates, err := api.getUpdates(ctx, 0)

				if err != nil {
					fmt.Printf("Got error: %v\n", err)
				}

				// Block current goroutine here until all updates will be readed by receivers:
				for i := range updates {
					res <- updates[i]
				}
			}
		}
	}()

	return res
}

func (api *TelegramAPI) sendMessage(ctx context.Context, req *requests.SendMessageRequest) error {
	serialized, err := req.Serialize()

	if err != nil {
		return fmt.Errorf("SERIALIZE_FAILED: %w", err)
	}

	buf := bytes.NewBuffer(serialized)

	res, err := api.request(ctx, sendMessageMethod, buf)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Got panic while getUpdates: %v\n", r)
		}

		if res != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		return err
	}

	type S map[string]interface{}

	r := new(S)

	b, _ := io.ReadAll(res.Body)

	json.Unmarshal(b, r)

	fmt.Printf("Response: %v+", *r)

	return nil
}

func (api *TelegramAPI) SendMessage(ctx context.Context, message *domain.MessageBuilder) error {
	return api.sendMessage(ctx, message.GetRequest())
}

func (api *TelegramAPI) Reply(ctx context.Context, update *responses.Update, text string) error {
	m := domain.NewMessageBuilder(strconv.Itoa(update.Message.Chat.Id), text)

	m.WithReplyToMessageId(update.Message.Id)

	return api.sendMessage(ctx, m.GetRequest())
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
