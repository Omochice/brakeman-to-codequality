package codequality_test

import (
	"bytes"
	"testing"

	"github.com/Omochice/brakeman-to-codequality/codequality"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	t.Run("writes valid JSON array", func(t *testing.T) {
		violations := []codequality.Violation{
			{
				Description: "Possible SQL injection",
				CheckName:   "SQL Injection",
				Fingerprint: "abc123",
				Severity:    "critical",
				Location: codequality.Location{
					Path:  "app/models/user.rb",
					Lines: codequality.Lines{Begin: 42},
				},
			},
		}

		var buf bytes.Buffer
		err := codequality.Write(violations, &buf)
		require.NoError(t, err)

		output := buf.String()
		require.Contains(t, output, "Possible SQL injection")
		require.Contains(t, output, "SQL Injection")
		require.Contains(t, output, "abc123")
		require.Contains(t, output, "critical")
	})

	t.Run("writes empty array", func(t *testing.T) {
		violations := []codequality.Violation{}

		var buf bytes.Buffer
		err := codequality.Write(violations, &buf)
		require.NoError(t, err)

		output := buf.String()
		require.Contains(t, output, "[")
		require.Contains(t, output, "]")
	})

	t.Run("output has no BOM", func(t *testing.T) {
		violations := []codequality.Violation{}

		var buf bytes.Buffer
		err := codequality.Write(violations, &buf)
		require.NoError(t, err)

		output := buf.Bytes()
		require.False(t, bytes.HasPrefix(output, []byte{0xEF, 0xBB, 0xBF}))
	})
}
