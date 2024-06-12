package errors

import "errors"

var ErrInvalidResponse = errors.New("INVALID_TELEGRAM_API_RESPONSE")

func NewErrInvalidResponse(message string) error {
	return newErrCustom(message, ErrInvalidResponse)
}
