package errors

import "fmt"

func newErrCustom(message string, cause error) error {
	return fmt.Errorf("%s: %w", message, cause)
}
