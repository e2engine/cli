package output_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ygrebnov/model/validation"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func TestNewSettings(t *testing.T) {
	got := output.NewSettings(
		"  json  ",
		true,
		true,
		true,
	)

	if got.Format != render.Format("json") {
		t.Errorf("Format = %q, want %q", got.Format, "json")
	}
	if !got.IsQuiet {
		t.Error("IsQuiet = false, want true")
	}
	if !got.IsVerbose {
		t.Error("IsVerbose = false, want true")
	}
	if !got.IsNoHeaders {
		t.Error("IsNoHeaders = false, want true")
	}
}

func TestValidationErrorOutput_TableHeaders(t *testing.T) {
	v := output.ValidationErrorOutput{}

	got := v.TableHeaders()
	want := []string{"FIELD", "RULE", "MESSAGE"}

	if len(got) != len(want) {
		t.Fatalf(
			"TableHeaders() length = %d, want %d",
			len(got),
			len(want),
		)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf(
				"TableHeaders()[%d] = %q, want %q",
				i,
				got[i],
				want[i],
			)
		}
	}
}

func TestValidationErrorOutput_TableRows(t *testing.T) {
	v := output.ValidationErrorOutput{
		Field:   "Name",
		Rule:    "min",
		Message: "name is too short",
	}

	got := v.TableRows()

	if len(got) != 1 {
		t.Fatalf("TableRows() length = %d, want 1", len(got))
	}

	want := []string{
		"Name",
		"min",
		"name is too short",
	}

	if len(got[0]) != len(want) {
		t.Fatalf(
			"TableRows()[0] length = %d, want %d",
			len(got[0]),
			len(want),
		)
	}

	for i := range want {
		if got[0][i] != want[i] {
			t.Errorf(
				"TableRows()[0][%d] = %q, want %q",
				i,
				got[0][i],
				want[i],
			)
		}
	}
}

func TestRenderValidationError_Nil(t *testing.T) {
	var buf bytes.Buffer

	err := output.RenderValidationError(
		&buf,
		nil,
		output.Settings{
			Format: render.Format("json"),
		},
	)
	if err != nil {
		t.Fatalf("RenderValidationError() error = %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf(
			"RenderValidationError() output = %q, want empty",
			buf.String(),
		)
	}
}

func TestRenderValidationError_Empty(t *testing.T) {
	var buf bytes.Buffer

	err := output.RenderValidationError(
		&buf,
		&validation.Error{},
		output.Settings{
			Format: render.Format("json"),
		},
	)
	if err != nil {
		t.Fatalf("RenderValidationError() error = %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf(
			"RenderValidationError() output = %q, want empty",
			buf.String(),
		)
	}
}

func TestRenderValidationError_Single(t *testing.T) {
	validationErr := &validation.Error{}
	validationErr.Addf(
		"Name",
		"min",
		errors.New("name is too short"),
	)

	var buf bytes.Buffer

	err := output.RenderValidationError(
		&buf,
		validationErr,
		output.Settings{
			Format: render.Format("json"),
		},
	)
	if err != nil {
		t.Fatalf("RenderValidationError() error = %v", err)
	}

	var got output.ValidationErrorOutput
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	want := output.ValidationErrorOutput{
		Valid:   false,
		Field:   "Name",
		Rule:    "min",
		Message: "name is too short",
	}

	if got != want {
		t.Errorf("RenderValidationError() = %#v, want %#v", got, want)
	}
}

func TestRenderValidationError_PreservesFieldOrder(t *testing.T) {
	validationErr := &validation.Error{}

	validationErr.Addf(
		"Version",
		"required",
		errors.New("version is required"),
	)
	validationErr.Addf(
		"Name",
		"min",
		errors.New("name is too short"),
	)
	validationErr.Addf(
		"Version",
		"format",
		errors.New("version has invalid format"),
	)

	var buf bytes.Buffer

	err := output.RenderValidationError(
		&buf,
		validationErr,
		output.Settings{
			Format: render.Format("json"),
		},
	)
	if err != nil {
		t.Fatalf("RenderValidationError() error = %v", err)
	}

	got := buf.String()

	versionRequired := strings.Index(
		got,
		"version is required",
	)
	versionFormat := strings.Index(
		got,
		"version has invalid format",
	)
	name := strings.Index(
		got,
		"name is too short",
	)

	if versionRequired < 0 ||
		versionFormat < 0 ||
		name < 0 {
		t.Fatalf(
			"RenderValidationError() output = %q, expected all validation errors",
			got,
		)
	}

	// Fields retain first-occurrence order and issues belonging
	// to the same field remain grouped together.
	if !(versionRequired < versionFormat &&
		versionFormat < name) {
		t.Errorf(
			"RenderValidationError() output order = %q, want Version issues before Name issue",
			got,
		)
	}
}

func TestRenderValidationError_NoHeaders(t *testing.T) {
	validationErr := &validation.Error{}
	validationErr.Addf(
		"Name",
		"required",
		errors.New("name is required"),
	)

	var buf bytes.Buffer

	err := output.RenderValidationError(
		&buf,
		validationErr,
		output.Settings{
			Format:      render.Format("table"),
			IsNoHeaders: true,
		},
	)
	if err != nil {
		t.Fatalf("RenderValidationError() error = %v", err)
	}

	got := buf.String()

	if strings.Contains(got, "FIELD") {
		t.Errorf(
			"RenderValidationError() output = %q, want no headers",
			got,
		)
	}

	if !strings.Contains(got, "Name") {
		t.Errorf(
			"RenderValidationError() output = %q, want Name",
			got,
		)
	}
}
