package errors

import (
	"errors"
)

var ErrValidate = errors.New("err while validate")

func NewErrValidate(message string) error {
	return newErrCustom(message, ErrValidate)
}
