package errors_test

import (
	"errors"
	"testing"

	internalerrors "github.com/e2engine/cli/internal/errors"
)

func TestUserError(t *testing.T) {
	cause := errors.New("underlying error")

	err := internalerrors.NewUserError(
		"operation failed",
		cause,
	)

	if got, want := err.Error(), "operation failed"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	if !errors.Is(err, cause) {
		t.Errorf("errors.Is() = false, want true")
	}
}

func TestUserError_NilCause(t *testing.T) {
	err := internalerrors.NewUserError(
		"operation failed",
		nil,
	)

	if got, want := err.Error(), "operation failed"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	if got := errors.Unwrap(err); got != nil {
		t.Errorf("errors.Unwrap() = %v, want nil", got)
	}
}
