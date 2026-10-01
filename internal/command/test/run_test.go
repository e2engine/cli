package test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestRunHandler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		RunTest(
			gomock.Any(),
			"environment-ref",
			"test-ref",
		).
		Return("execution-id", nil)

	var writer bytes.Buffer

	err := NewRunHandler(service).
		Run(context.Background(), RunParams{
			EnvironmentReference: "environment-ref",
			TestReference:        "test-ref",
			Writer:               &writer,
		})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := writer.String(); !strings.Contains(
		got,
		"created test execution with id: execution-id",
	) {
		t.Errorf(
			"output = %q, want execution creation message",
			got,
		)
	}
}

func TestRunHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		RunTest(
			gomock.Any(),
			"environment-ref",
			"test-ref",
		).
		Return("execution-id", nil)

	var writer bytes.Buffer

	params := RunParams{
		EnvironmentReference: "environment-ref",
		TestReference:        "test-ref",
		Writer:               &writer,
	}
	params.Output.IsQuiet = true

	err := NewRunHandler(service).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := writer.String(), "execution-id\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestRunHandler_Run_Verbose(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		RunTest(
			gomock.Any(),
			"environment-ref",
			"test-ref",
		).
		Return("execution-id", nil)

	var writer bytes.Buffer

	params := RunParams{
		EnvironmentReference: "environment-ref",
		TestReference:        "test-ref",
		Writer:               &writer,
	}
	params.Output.IsVerbose = true

	err := NewRunHandler(service).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := writer.String(),
		"created test execution with id: execution-id\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestRunHandler_Run_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	wantErr := errors.New("run failed")

	service.EXPECT().
		RunTest(
			gomock.Any(),
			"environment-ref",
			"test-ref",
		).
		Return("", wantErr)

	var writer bytes.Buffer

	err := NewRunHandler(service).
		Run(context.Background(), RunParams{
			EnvironmentReference: "environment-ref",
			TestReference:        "test-ref",
			Writer:               &writer,
		})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}

	if writer.Len() != 0 {
		t.Errorf("output = %q, want empty", writer.String())
	}
}
