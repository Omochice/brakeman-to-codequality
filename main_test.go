package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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
	t.Run("reads from stdin when source is dash", func(t *testing.T) {
		input := `{"warnings":[{"warning_type":"SQL Injection","message":"Possible SQL injection","file":"app/models/user.rb","line":42,"confidence":"High","code":"User.where(...)"}]}`

		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(input),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command([]string{"-"}, inout)
		if exitCode != 0 {
			t.Fatalf("got %v, want %v\nstderr: %s", exitCode, 0, stderr.String())
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

	t.Run("reads from file when source is a file path", func(t *testing.T) {
		content := `{"warnings":[{"warning_type":"XSS","message":"Cross-site scripting","file":"app/views/index.erb","line":10,"confidence":"Medium"}]}`
		dir := t.TempDir()
		path := filepath.Join(dir, "report.json")
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(""),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command([]string{path}, inout)
		if exitCode != 0 {
			t.Fatalf("got %v, want %v\nstderr: %s", exitCode, 0, stderr.String())
		}

		output := stdout.String()
		if !strings.Contains(output, "Cross-site scripting") {
			t.Fatalf("expected %q to contain %q", output, "Cross-site scripting")
		}
	})

	t.Run("handles empty warnings with dash", func(t *testing.T) {
		input := `{"warnings":[]}`

		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(input),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command([]string{"-"}, inout)
		if exitCode != 0 {
			t.Fatalf("got %v, want %v", exitCode, 0)
		}
		if stderr.String() != "" {
			t.Fatalf("expected empty stderr, got %q", stderr.String())
		}

		var result []any
		if err := json.NewDecoder(&stdout).Decode(&result); err != nil {
			t.Fatalf("failed to decode output as JSON: %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected empty array, got %d elements", len(result))
		}
	})

	t.Run("returns non-zero exit code for invalid JSON from stdin", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(`{invalid json`),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command([]string{"-"}, inout)
		if exitCode != 1 {
			t.Fatalf("got %v, want %v", exitCode, 1)
		}
		if !strings.Contains(stderr.String(), "Error:") {
			t.Fatalf("expected %q to contain %q", stderr.String(), "Error:")
		}
	})

	t.Run("returns non-zero exit code when no arguments given", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(""),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command([]string{}, inout)
		if exitCode != 1 {
			t.Fatalf("got %v, want %v", exitCode, 1)
		}
	})

	t.Run("returns non-zero exit code when multiple arguments given", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(""),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		exitCode := command([]string{"a.json", "b.json"}, inout)
		if exitCode != 1 {
			t.Fatalf("got %v, want %v", exitCode, 1)
		}
	})

	t.Run("returns non-zero exit code when file does not exist", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		inout := &cli.ProcInout{
			Stdin:  strings.NewReader(""),
			Stdout: &stdout,
			Stderr: &stderr,
		}

		missingFile := filepath.Join(t.TempDir(), "file.json")
		exitCode := command([]string{missingFile}, inout)
		if exitCode != 1 {
			t.Fatalf("got %v, want %v", exitCode, 1)
		}
	})
}
