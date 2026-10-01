package environment

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
		DeleteEnvironment(gomock.Any(), "environment-ref").
		Return(&model.Environment{
			ID:   "environment-id",
			Name: "environment-name",
		}, nil)

	handler := NewDeleteHandler(service)

	var writer bytes.Buffer

	err := handler.Run(context.Background(), DeleteParams{
		Reference: "environment-ref",
		Writer:    &writer,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := writer.String()

	if !strings.Contains(got, "deleted environment") {
		t.Errorf("output = %q, want deletion message", got)
	}
	if !strings.Contains(got, "environment-name") {
		t.Errorf("output = %q, want environment name", got)
	}
	if !strings.Contains(got, "environment-id") {
		t.Errorf("output = %q, want environment ID", got)
	}
}

func TestDeleteHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		DeleteEnvironment(gomock.Any(), "environment-ref").
		Return(&model.Environment{
			ID:   "environment-id",
			Name: "environment-name",
		}, nil)

	handler := NewDeleteHandler(service)

	var writer bytes.Buffer

	params := DeleteParams{
		Reference: "environment-ref",
		Writer:    &writer,
	}
	params.Output.IsQuiet = true

	err := handler.Run(context.Background(), params)
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
		DeleteEnvironment(gomock.Any(), "environment-ref").
		Return(nil, wantErr)

	handler := NewDeleteHandler(service)

	err := handler.Run(context.Background(), DeleteParams{
		Reference: "environment-ref",
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}
