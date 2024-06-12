package requests

import (
	"encoding/json"

	"github.com/EXPORTER-DEV/go-telegram-bot/pkg/bot/cmd/api/definitions"
)

type GetUpdatesRequest struct {
	Offset           int             `json:"offset"`
	Limit            int             `json:"limit"`
	TimeoutInSeconds float64         `json:"timeout"`
	AllowedUpdates   []AllowedUpdate `json:"allowed_updates,omitempty"`
}

type AllowedUpdate string

var MessageUpdate AllowedUpdate = "message"
var CallbackQueryUpdate AllowedUpdate = "callback_query"

func (req *GetUpdatesRequest) Serialize() ([]byte, error) {
	return json.Marshal(req)
}

func (req *GetUpdatesRequest) Validate() error {
	return nil
}

func NewGetUpdatesRequest(
	offset int,
	limit int,
	timeoutInSeconds float64,
	allowedUpdates []AllowedUpdate,
) definitions.Requester {
	return &GetUpdatesRequest{
		offset,
		limit,
		timeoutInSeconds,
		allowedUpdates,
	}
}
