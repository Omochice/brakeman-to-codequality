package codequality_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Omochice/brakeman-to-codequality/codequality"
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
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "Possible SQL injection") {
			t.Fatalf("expected %q to contain %q", output, "Possible SQL injection")
		}
		if !strings.Contains(output, "SQL Injection") {
			t.Fatalf("expected %q to contain %q", output, "SQL Injection")
		}
		if !strings.Contains(output, "abc123") {
			t.Fatalf("expected %q to contain %q", output, "abc123")
		}
		if !strings.Contains(output, "critical") {
			t.Fatalf("expected %q to contain %q", output, "critical")
		}
	})

	t.Run("writes empty array", func(t *testing.T) {
		violations := []codequality.Violation{}

		var buf bytes.Buffer
		err := codequality.Write(violations, &buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "[") {
			t.Fatalf("expected %q to contain %q", output, "[")
		}
		if !strings.Contains(output, "]") {
			t.Fatalf("expected %q to contain %q", output, "]")
		}
	})

	t.Run("output has no BOM", func(t *testing.T) {
		violations := []codequality.Violation{}

		var buf bytes.Buffer
		err := codequality.Write(violations, &buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := buf.Bytes()
		if bytes.HasPrefix(output, []byte{0xEF, 0xBB, 0xBF}) {
			t.Fatalf("expected false, got true")
		}
	})
}
