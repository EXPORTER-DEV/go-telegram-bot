package errors

import "errors"

func IsErrCausedBy(err error, target error) bool {
	return errors.Is(err, target)
}
