package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Omochice/brakeman-to-codequality/cli"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	t.Run("writes error to writer and returns 1", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode := handleError(&buf, errors.New("test error"))

		require.Equal(t, 1, exitCode)
		require.Contains(t, buf.String(), "Error:")
		require.Contains(t, buf.String(), "test error")
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
		require.Equal(t, 0, exitCode)
		require.Empty(t, stderr.String())

		output := stdout.String()
		require.Contains(t, output, "Possible SQL injection")
		require.Contains(t, output, "SQL Injection")
		require.Contains(t, output, "critical")
		require.Contains(t, output, "app/models/user.rb")
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
		require.Equal(t, 0, exitCode)
		require.Empty(t, stderr.String())

		output := stdout.String()
		require.Contains(t, output, "[")
		require.Contains(t, output, "]")
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
		require.Equal(t, 1, exitCode)
		require.Contains(t, stderr.String(), "Error:")
	})
}
