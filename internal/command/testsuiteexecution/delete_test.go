package testsuiteexecution

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/e2engine/core/model"
	"go.uber.org/mock/gomock"
)

func TestDeleteHandler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		DeleteTestSuiteExecution(gomock.Any(), "execution-ref").
		Return(&model.TestSuiteExecution{
			ID: "execution-id",
		}, nil)

	var writer bytes.Buffer

	err := NewDeleteHandler(service).
		Run(context.Background(), DeleteParams{
			Reference: "execution-ref",
			Writer:    &writer,
		})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := writer.String(),
		"deleted testsuite execution with id: execution-id\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestDeleteHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		DeleteTestSuiteExecution(gomock.Any(), "execution-ref").
		Return(&model.TestSuiteExecution{
			ID: "execution-id",
		}, nil)

	var writer bytes.Buffer

	params := DeleteParams{
		Reference: "execution-ref",
		Writer:    &writer,
	}
	params.Output.IsQuiet = true

	err := NewDeleteHandler(service).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if writer.Len() != 0 {
		t.Errorf("output = %q, want empty", writer.String())
	}
}

func TestDeleteHandler_Run_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	wantErr := errors.New("delete failed")

	service.EXPECT().
		DeleteTestSuiteExecution(gomock.Any(), "execution-ref").
		Return(nil, wantErr)

	var writer bytes.Buffer

	err := NewDeleteHandler(service).
		Run(context.Background(), DeleteParams{
			Reference: "execution-ref",
			Writer:    &writer,
		})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}

	if writer.Len() != 0 {
		t.Errorf("output = %q, want empty", writer.String())
	}
}
