package errors

import (
	"errors"
)

var ErrSerialize = errors.New("err serialize")

func NewErrSerialize(message string) error {
	return newErrCustom(message, ErrSerialize)
}
