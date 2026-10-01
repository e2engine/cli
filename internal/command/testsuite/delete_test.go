package testsuite

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
		DeleteTestSuite(gomock.Any(), "testsuite-ref").
		Return(&model.TestSuite{
			ID:   "testsuite-id",
			Name: "testsuite-name",
		}, nil)

	var writer bytes.Buffer

	err := NewDeleteHandler(service).
		Run(context.Background(), DeleteParams{
			Reference: "testsuite-ref",
			Writer:    &writer,
		})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := writer.String()

	if !strings.Contains(got, "deleted testsuite") {
		t.Errorf("output = %q, want deletion message", got)
	}
	if !strings.Contains(got, "testsuite-name") {
		t.Errorf("output = %q, want testsuite name", got)
	}
	if !strings.Contains(got, "testsuite-id") {
		t.Errorf("output = %q, want testsuite ID", got)
	}
}

func TestDeleteHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		DeleteTestSuite(gomock.Any(), "testsuite-ref").
		Return(&model.TestSuite{
			ID:   "testsuite-id",
			Name: "testsuite-name",
		}, nil)

	var writer bytes.Buffer

	params := DeleteParams{
		Reference: "testsuite-ref",
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
		DeleteTestSuite(gomock.Any(), "testsuite-ref").
		Return(nil, wantErr)

	err := NewDeleteHandler(service).
		Run(context.Background(), DeleteParams{
			Reference: "testsuite-ref",
		})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}
