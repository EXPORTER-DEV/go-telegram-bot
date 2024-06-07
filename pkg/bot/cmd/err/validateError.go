package err

import (
	"errors"
	"fmt"
)

var ErrValidate = errors.New("err while validate")

func NewValidateError(message string) error {
	return fmt.Errorf(message+": %w", ErrValidate)
}
