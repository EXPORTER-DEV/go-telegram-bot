package api

import (
	"context"
	"testing"
	"time"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/domain/builder"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/client"
	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/errors"
)

func getMockClientWithRequestor(t *testing.T) (*client.MockClient, Requester, error) {
	mockClient := client.NewMockClient(t)
	api, err := New("", 5, "", time.Duration(time.Second*60), mockClient)

	return mockClient, api, err
}

func getContext() context.Context {
	return context.Background()
}

func TestIncorrectSendMessageRequest(t *testing.T) {
	_, requestor, err := getMockClientWithRequestor(t)

	if err != nil {
		t.Fatalf("Got err while build api: %v", err)
	}

	ctx := getContext()

	err = requestor.SendMessage(ctx, builder.NewMessageBuilder("", ""))

	if !errors.IsErrCausedBy(err, errors.ErrValidate) {
		t.Fatalf("Got no error, expected: %+v", errors.ErrValidate)
	}
}
