package requests

import (
	"encoding/json"
)

type GetUpdatesRequest struct {
	Offset           int      `json:"offset"`
	Limit            int      `json:"limit"`
	TimeoutInSeconds float64  `json:"timeout"`
	AllowedUpdates   []string `json:"allowed_updates"`
}

func (req *GetUpdatesRequest) Serialize() ([]byte, error) {
	return json.Marshal(req)
}

func NewGetUpdatesRequest(offset int, limit int, timeoutInSeconds float64, allowedUpdates []string) *GetUpdatesRequest {
	return &GetUpdatesRequest{
		offset,
		limit,
		timeoutInSeconds,
		allowedUpdates,
	}
}

var GetUpdatesMethod = "getUpdates"
