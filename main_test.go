package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Omochice/brakeman-to-codequality/cli"
)

func TestHandleError(t *testing.T) {
	t.Run("writes error to writer and returns 1", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode := handleError(&buf, errors.New("test error"))

		if exitCode != 1 {
			t.Fatalf("got %v, want %v", exitCode, 1)
		}
		if !strings.Contains(buf.String(), "Error:") {
			t.Fatalf("expected %q to contain %q", buf.String(), "Error:")
		}
		if !strings.Contains(buf.String(), "test error") {
			t.Fatalf("expected %q to contain %q", buf.String(), "test error")
		}
	})
}

func TestEndToEnd(t *testing.T) {
	t.Run("complete conversion flow", func(t *testing.T) {
		input := `{"warnings":[{"warning_type":"SQL Injection","message":"Possible SQL injection","file":"app/models/user.rb","line":42,"confidence":"High","code":"User.where(...)"}]}`

		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(input),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command(nil, inout)
		if exitCode != 0 {
			t.Fatalf("got %v, want %v", exitCode, 0)
		}
		if stderr.String() != "" {
			t.Fatalf("expected empty string, got %q", stderr.String())
		}

		output := stdout.String()
		if !strings.Contains(output, "Possible SQL injection") {
			t.Fatalf("expected %q to contain %q", output, "Possible SQL injection")
		}
		if !strings.Contains(output, "SQL Injection") {
			t.Fatalf("expected %q to contain %q", output, "SQL Injection")
		}
		if !strings.Contains(output, "critical") {
			t.Fatalf("expected %q to contain %q", output, "critical")
		}
		if !strings.Contains(output, "app/models/user.rb") {
			t.Fatalf("expected %q to contain %q", output, "app/models/user.rb")
		}
	})

	t.Run("handles empty warnings", func(t *testing.T) {
		input := `{"warnings":[]}`

		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(input),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command(nil, inout)
		if exitCode != 0 {
			t.Fatalf("got %v, want %v", exitCode, 0)
		}
		if stderr.String() != "" {
			t.Fatalf("expected empty string, got %q", stderr.String())
		}

		output := stdout.String()
		if !strings.Contains(output, "[") {
			t.Fatalf("expected %q to contain %q", output, "[")
		}
		if !strings.Contains(output, "]") {
			t.Fatalf("expected %q to contain %q", output, "]")
		}
	})

	t.Run("returns non-zero exit code for invalid JSON input", func(t *testing.T) {
		input := `{invalid json`

		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(input),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command(nil, inout)
		if exitCode != 1 {
			t.Fatalf("got %v, want %v", exitCode, 1)
		}
		if !strings.Contains(stderr.String(), "Error:") {
			t.Fatalf("expected %q to contain %q", stderr.String(), "Error:")
		}
	})
}
