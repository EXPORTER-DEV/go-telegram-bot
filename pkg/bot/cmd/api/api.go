package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/domain/builder"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/requests"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/responses"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/client"
)

var ErrInvalidResponse = errors.New("INVALID_TELEGRAM_API_RESPONSE")
var ErrInvalidArgument = errors.New("INVALID_ARGUMENT")

type Requester interface {
	Poll(ctx context.Context) chan *responses.Update
	SendMessage(ctx context.Context, message builder.MessageBuilder) error
	ReplyTo(ctx context.Context, target *responses.Update, text string) error
	SetDebugMode(debug bool)
}

type API struct {
	token      string
	offset     int
	Limit      int
	URL        *url.URL
	Timeout    time.Duration
	client     client.Client
	maxRetries int
	debug      bool
}

type method string

var getUpdatesMethod method = "getUpdates"
var sendMessageMethod method = "sendMessage"

func (api *API) request(ctx context.Context, method method, body io.Reader) (*http.Response, error) {
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

func (api *API) getUpdates(ctx context.Context, retry int) ([]*responses.Update, error) {
	req := requests.NewGetUpdatesRequest(
		api.offset,
		api.Limit,
		api.Timeout.Seconds(),
		[]requests.AllowedUpdate{requests.MessageUpdate, requests.CallbackQueryUpdate},
	)

	serialized, err := req.Serialize()

	if err != nil {
		return nil, fmt.Errorf("SERIALIZE_FAILED: %w", err)
	}

	buf := bytes.NewBuffer(serialized)

	res, err := api.request(ctx, getUpdatesMethod, buf)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Got panic while getUpdates: %v\n", r)
		}

		if res != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		if os.IsTimeout(err) && api.maxRetries > retry {
			retry += 1
			log.Printf("Timeouted, going for retry: %v/%v\n", retry, api.maxRetries)
			return api.getUpdates(ctx, retry)
		}

		return nil, err
	}

	decoder := json.NewDecoder(res.Body)

	if api.debug {
		bodyReader, _ := res.Request.GetBody()

		requestBody, _ := io.ReadAll(bodyReader)

		log.Printf("Request body: %s\n", requestBody)
	}

	var response = new(responses.GetUpdatesResponse)

	err = decoder.Decode(response)

	if err != nil {
		return nil, err
	}

	if !response.Ok {
		log.Printf("Got incorrect response from Telegram API: %+v", response)
		return nil, fmt.Errorf("GOT_INCORRECT_RESPONSE: %w", ErrInvalidResponse)
	}

	log.Printf("Response: %+v\n", response)

	if len(response.Result) > 0 {
		// Cause result is in chronic order, so last has the highest updateId:
		api.offset = response.Result[len(response.Result)-1].UpdateId + 1
		log.Printf("Updated offset: %+v\n", api.offset)
	}

	return response.Result, nil
}

func (api *API) Poll(ctx context.Context) chan *responses.Update {
	res := make(chan *responses.Update)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Canceled polling cause context done\n")
				// Probably extra operation, but let it exist here:
				close(res)
				return
			default:
				updates, err := api.getUpdates(ctx, 0)

				if err != nil {
					log.Printf("Got error: %v\n", err)
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

func (api *API) sendMessage(ctx context.Context, req *requests.SendMessageRequest) error {
	if err := req.Validate(); err != nil {
		log.Printf("Got invalid message: %v, req: %+v", err, req)

		return err
	}

	serialized, err := req.Serialize()

	if err != nil {
		return fmt.Errorf("SERIALIZE_FAILED: %w", err)
	}

	buf := bytes.NewBuffer(serialized)

	res, err := api.request(ctx, sendMessageMethod, buf)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Got panic while getUpdates: %v\n", r)
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

	if api.debug {
		log.Printf("Response: %v+", *r)
	}

	return nil
}

func (api *API) SendMessage(ctx context.Context, message builder.MessageBuilder) error {
	return api.sendMessage(ctx, message.GetRequest())
}

func (api *API) ReplyTo(ctx context.Context, update *responses.Update, text string) error {
	m := builder.NewMessageBuilder(strconv.Itoa(update.Message.Chat.Id), text)

	m.WithReplyToMessageId(update.Message.Id)

	return api.sendMessage(ctx, m.GetRequest())
}

func (api *API) SetDebugMode(debug bool) {
	api.debug = debug
}

func New(token string, limit int, rawURL string, timeout time.Duration, client client.Client) (Requester, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &API{
		token:      token,
		offset:     0,
		Limit:      limit,
		URL:        url,
		Timeout:    timeout,
		client:     client,
		maxRetries: 5,
		debug:      false,
	}, nil
}
