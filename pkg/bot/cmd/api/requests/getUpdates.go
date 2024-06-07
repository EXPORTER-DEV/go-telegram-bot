package requests

import (
	"encoding/json"
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

func NewGetUpdatesRequest(offset int, limit int, timeoutInSeconds float64, allowedUpdates []AllowedUpdate) *GetUpdatesRequest {
	return &GetUpdatesRequest{
		offset,
		limit,
		timeoutInSeconds,
		allowedUpdates,
	}
}
