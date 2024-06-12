package errors

import (
	"errors"
	"testing"
)

var err = errors.New("TEST")

func TestIsErrCausedBy(t *testing.T) {
	errInvalidResponse := newErrCustom("TEST MESSAGE", err)

	if !IsErrCausedBy(errInvalidResponse, err) {
		t.Fatalf("Error: %+v not match: %+v", errInvalidResponse, ErrInvalidResponse)
	}
}
