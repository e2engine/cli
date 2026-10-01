package service

import (
	"errors"
	"testing"

	coreerrors "github.com/e2engine/core/pkg/errors"
	"github.com/ygrebnov/model/validation"

	internalerrors "github.com/e2engine/cli/internal/errors"
)

func TestService_Close_NilCloser(t *testing.T) {
	service := &Service{}

	if err := service.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestService_Close(t *testing.T) {
	var closed bool

	service := &Service{
		closer: newCloser(
			&serviceTestCloser{
				close: func() error {
					closed = true
					return nil
				},
			},
		),
	}

	if err := service.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	if !closed {
		t.Fatal("Close() did not close resource")
	}
}

func TestCloserStack_Close_ReverseOrder(t *testing.T) {
	var order []int

	closer := newCloser(
		&serviceTestCloser{
			close: func() error {
				order = append(order, 1)
				return nil
			},
		},
		&serviceTestCloser{
			close: func() error {
				order = append(order, 2)
				return nil
			},
		},
		&serviceTestCloser{
			close: func() error {
				order = append(order, 3)
				return nil
			},
		},
	)

	if err := closer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	want := []int{3, 2, 1}

	if len(order) != len(want) {
		t.Fatalf(
			"close order = %v, want %v",
			order,
			want,
		)
	}

	for i := range want {
		if order[i] != want[i] {
			t.Fatalf(
				"close order = %v, want %v",
				order,
				want,
			)
		}
	}
}

func TestCloserStack_Close_AllResourcesOnError(t *testing.T) {
	err1 := errors.New("first")
	err2 := errors.New("second")

	var order []int

	closer := newCloser(
		&serviceTestCloser{
			close: func() error {
				order = append(order, 1)
				return err1
			},
		},
		&serviceTestCloser{
			close: func() error {
				order = append(order, 2)
				return nil
			},
		},
		&serviceTestCloser{
			close: func() error {
				order = append(order, 3)
				return err2
			},
		},
	)

	err := closer.Close()
	if err == nil {
		t.Fatal("Close() error = nil, want error")
	}

	if !errors.Is(err, err1) {
		t.Fatalf(
			"Close() error = %v, want errors.Is(..., err1)",
			err,
		)
	}

	if !errors.Is(err, err2) {
		t.Fatalf(
			"Close() error = %v, want errors.Is(..., err2)",
			err,
		)
	}

	wantOrder := []int{3, 2, 1}

	if len(order) != len(wantOrder) {
		t.Fatalf(
			"close order = %v, want %v",
			order,
			wantOrder,
		)
	}

	for i := range wantOrder {
		if order[i] != wantOrder[i] {
			t.Fatalf(
				"close order = %v, want %v",
				order,
				wantOrder,
			)
		}
	}
}

func TestCloserStack_Add(t *testing.T) {
	var order []int

	closer := newCloser(
		&serviceTestCloser{
			close: func() error {
				order = append(order, 1)
				return nil
			},
		},
	)

	closer.Add(
		&serviceTestCloser{
			close: func() error {
				order = append(order, 2)
				return nil
			},
		},
	)

	if err := closer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	want := []int{2, 1}

	for i := range want {
		if order[i] != want[i] {
			t.Fatalf(
				"close order = %v, want %v",
				order,
				want,
			)
		}
	}
}

func TestUserFacingValidationError_NotInvalidSpec(t *testing.T) {
	original := errors.New("test error")

	actual := userFacingValidationError(
		original,
		"environment",
	)

	if actual != original {
		t.Fatalf(
			"userFacingValidationError() = %v, want original error %v",
			actual,
			original,
		)
	}
}

func TestUserFacingValidationError_InvalidSpecWithoutValidationError(t *testing.T) {
	actual := userFacingValidationError(
		coreerrors.ErrInvalidSpec,
		"environment",
	)

	if actual != coreerrors.ErrInvalidSpec {
		t.Fatalf(
			"userFacingValidationError() = %v, want original error %v",
			actual,
			coreerrors.ErrInvalidSpec,
		)
	}
}

func TestUserFacingValidationError(t *testing.T) {
	validationErr := &validation.Error{}
	validationErr.Addf(
		"Name",
		"required",
		errors.New("must not be empty"),
	)
	validationErr.Addf(
		"Services.HTTP.Port",
		"required",
		errors.New("must not be empty"),
	)
	validationErr.Addf(
		"Name",
		"another-rule",
		errors.New("another error"),
	)

	err := errors.Join(
		coreerrors.ErrInvalidSpec,
		validationErr,
	)

	actual := userFacingValidationError(
		err,
		"environment",
	)

	var userErr *internalerrors.UserError
	if !errors.As(actual, &userErr) {
		t.Fatalf(
			"userFacingValidationError() error = %T, want *errors.UserError",
			actual,
		)
	}

	wantMessage := "invalid environment spec, fields: Name, Services.HTTP.Port"
	if userErr.Message != wantMessage {
		t.Errorf(
			"UserError.Message = %q, want %q",
			userErr.Message,
			wantMessage,
		)
	}

	if userErr.Err != err {
		t.Errorf(
			"UserError.Err = %v, want original error %v",
			userErr.Err,
			err,
		)
	}

	if !errors.Is(actual, coreerrors.ErrInvalidSpec) {
		t.Errorf(
			"userFacingValidationError() does not wrap ErrInvalidSpec",
		)
	}

	if !errors.Is(actual, validationErr) {
		t.Errorf(
			"userFacingValidationError() does not wrap validation error",
		)
	}
}

type serviceTestCloser struct {
	close func() error
}

func (c *serviceTestCloser) Close() error {
	return c.close()
}
