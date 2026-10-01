package testsuite

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/e2engine/core/model"
	"github.com/ygrebnov/render"
	"go.uber.org/mock/gomock"

	clierrors "github.com/e2engine/cli/pkg/errors"
)

func TestCreateHandler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateTestSuite(gomock.Any(), "testsuite.yaml").
		Return(&model.TestSuite{
			ID:   "testsuite-id",
			Name: "testsuite-name",
		}, nil)

	var writer bytes.Buffer

	err := NewCreateHandler(service).
		Run(context.Background(), CreateParams{
			SpecPath: "testsuite.yaml",
			Writer:   &writer,
		})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := writer.String()

	if !strings.Contains(got, "created testsuite") {
		t.Errorf("output = %q, want creation message", got)
	}
	if !strings.Contains(got, "testsuite-name") {
		t.Errorf("output = %q, want testsuite name", got)
	}
	if !strings.Contains(got, "testsuite-id") {
		t.Errorf("output = %q, want testsuite ID", got)
	}
}

func TestCreateHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateTestSuite(gomock.Any(), "testsuite.yaml").
		Return(&model.TestSuite{
			ID:   "testsuite-id",
			Name: "testsuite-name",
		}, nil)

	var writer bytes.Buffer

	params := CreateParams{
		SpecPath: "testsuite.yaml",
		Writer:   &writer,
	}
	params.Output.IsQuiet = true

	err := NewCreateHandler(service).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := writer.String(), "testsuite-id\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestCreateHandler_Run_Verbose(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateTestSuite(gomock.Any(), "testsuite.yaml").
		Return(&model.TestSuite{
			ID:          "testsuite-id",
			Name:        "testsuite-name",
			Description: "testsuite description",
		}, nil)

	var writer bytes.Buffer

	params := CreateParams{
		SpecPath: "testsuite.yaml",
		Writer:   &writer,
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatYAML

	err := NewCreateHandler(service).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(writer.String(), "testsuite-name") {
		t.Errorf(
			"output does not contain testsuite name: %q",
			writer.String(),
		)
	}

	if !strings.Contains(writer.String(), "testsuite description") {
		t.Errorf(
			"output does not contain description: %q",
			writer.String(),
		)
	}
}

func TestCreateHandler_Run_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	wantErr := errors.New("create failed")

	service.EXPECT().
		CreateTestSuite(gomock.Any(), "testsuite.yaml").
		Return(nil, wantErr)

	err := NewCreateHandler(service).
		Run(context.Background(), CreateParams{
			SpecPath: "testsuite.yaml",
		})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}

func TestCreateHandler_Run_RenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateTestSuite(gomock.Any(), "testsuite.yaml").
		Return(&model.TestSuite{Name: "testsuite"}, nil)

	params := CreateParams{
		SpecPath: "testsuite.yaml",
		Writer: failingWriter{
			err: errors.New("write failed"),
		},
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatYAML

	err := NewCreateHandler(service).
		Run(context.Background(), params)

	if !errors.Is(err, clierrors.ErrCannotRenderOutput) {
		t.Fatalf(
			"Run() error = %v, want ErrCannotRenderOutput",
			err,
		)
	}
}
