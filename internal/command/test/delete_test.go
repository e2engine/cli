package test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/e2engine/core/model"
	"go.uber.org/mock/gomock"
)

func TestDeleteHandler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		DeleteTest(gomock.Any(), "test-ref").
		Return(&model.Test{
			ID:   "test-id",
			Name: "test-name",
		}, nil)

	var writer bytes.Buffer

	err := NewDeleteHandler(service).
		Run(context.Background(), DeleteParams{
			Reference: "test-ref",
			Writer:    &writer,
		})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := writer.String()

	if !strings.Contains(got, "deleted test") {
		t.Errorf("output = %q, want deletion message", got)
	}
	if !strings.Contains(got, "test-name") {
		t.Errorf("output = %q, want test name", got)
	}
	if !strings.Contains(got, "test-id") {
		t.Errorf("output = %q, want test ID", got)
	}
}

func TestDeleteHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		DeleteTest(gomock.Any(), "test-ref").
		Return(&model.Test{
			ID:   "test-id",
			Name: "test-name",
		}, nil)

	var writer bytes.Buffer

	params := DeleteParams{
		Reference: "test-ref",
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
		DeleteTest(gomock.Any(), "test-ref").
		Return(nil, wantErr)

	err := NewDeleteHandler(service).
		Run(context.Background(), DeleteParams{
			Reference: "test-ref",
		})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}
