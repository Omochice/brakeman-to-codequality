package main

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
	"github.com/Omochice/brakeman-to-codequality/cli"
	"github.com/stretchr/testify/require"
)

func TestMapSeverity(t *testing.T) {
	tests := []struct {
		name       string
		confidence string
		want       string
	}{
		{
			name:       "High confidence maps to critical",
			confidence: "High",
			want:       "critical",
		},
		{
			name:       "high confidence (lowercase) maps to critical",
			confidence: "high",
			want:       "critical",
		},
		{
			name:       "HIGH confidence (uppercase) maps to critical",
			confidence: "HIGH",
			want:       "critical",
		},
		{
			name:       "Medium confidence maps to major",
			confidence: "Medium",
			want:       "major",
		},
		{
			name:       "medium confidence (lowercase) maps to major",
			confidence: "medium",
			want:       "major",
		},
		{
			name:       "Weak confidence maps to minor",
			confidence: "Weak",
			want:       "minor",
		},
		{
			name:       "weak confidence (lowercase) maps to minor",
			confidence: "weak",
			want:       "minor",
		},
		{
			name:       "Low confidence maps to minor",
			confidence: "Low",
			want:       "minor",
		},
		{
			name:       "low confidence (lowercase) maps to minor",
			confidence: "low",
			want:       "minor",
		},
		{
			name:       "Unknown confidence maps to info",
			confidence: "Unknown",
			want:       "info",
		},
		{
			name:       "Empty confidence maps to info",
			confidence: "",
			want:       "info",
		},
		{
			name:       "Invalid confidence maps to info",
			confidence: "InvalidValue",
			want:       "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapSeverity(tt.confidence)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateFingerprint(t *testing.T) {
	t.Run("generates consistent fingerprint for same input", func(t *testing.T) {
		fp1 := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "User.where(...)")
		fp2 := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "User.where(...)")
		require.Equal(t, fp1, fp2)
	})

	t.Run("generates different fingerprints for different inputs", func(t *testing.T) {
		fp1 := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := GenerateFingerprint("app/models/user.rb", 43, "SQL Injection", "Possible SQL injection", "")
		require.NotEqual(t, fp1, fp2)
	})

	t.Run("includes code field when present", func(t *testing.T) {
		fpWithCode := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "User.where(...)")
		fpWithoutCode := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		require.NotEqual(t, fpWithCode, fpWithoutCode)
	})

	t.Run("handles empty code field", func(t *testing.T) {
		fp := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		require.NotEmpty(t, fp)
	})

	t.Run("generates 64-character hex string", func(t *testing.T) {
		fp := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		require.Len(t, fp, 64)
		require.Regexp(t, regexp.MustCompile("^[0-9a-f]{64}$"), fp)
	})

	t.Run("different files produce different fingerprints", func(t *testing.T) {
		fp1 := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := GenerateFingerprint("app/models/post.rb", 42, "SQL Injection", "Possible SQL injection", "")
		require.NotEqual(t, fp1, fp2)
	})

	t.Run("different warning types produce different fingerprints", func(t *testing.T) {
		fp1 := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := GenerateFingerprint("app/models/user.rb", 42, "XSS", "Possible SQL injection", "")
		require.NotEqual(t, fp1, fp2)
	})

	t.Run("different messages produce different fingerprints", func(t *testing.T) {
		fp1 := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := GenerateFingerprint("app/models/user.rb", 42, "SQL Injection", "Confirmed SQL injection", "")
		require.NotEqual(t, fp1, fp2)
	})
}

func TestConvertWarnings(t *testing.T) {
	t.Run("converts valid warning correctly", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
				Code:        "User.where(...)",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 1)

		violation := violations[0]
		require.Equal(t, "Possible SQL injection", violation.Description)
		require.Equal(t, "SQL Injection", violation.CheckName)
		require.Equal(t, "critical", violation.Severity)
		require.Equal(t, "app/models/user.rb", violation.Location.Path)
		require.Equal(t, 42, violation.Location.Lines.Begin)
		require.NotEmpty(t, violation.Fingerprint)
	})

	t.Run("skips warning with missing file", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "",
				Line:        42,
				Confidence:  "High",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 0)
	})

	t.Run("skips warning with missing line", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        0,
				Confidence:  "High",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 0)
	})

	t.Run("skips warning with missing warning type", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 0)
	})

	t.Run("skips warning with missing message", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 0)
	})

	t.Run("removes ./ prefix from file path", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "./app/models/user.rb",
				Line:        42,
				Confidence:  "High",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 1)
		require.Equal(t, "app/models/user.rb", violations[0].Location.Path)
	})

	t.Run("handles empty array", func(t *testing.T) {
		warnings := []brakeman.Warning{}
		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 0)
	})

	t.Run("processes multiple warnings", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
			},
			{
				WarningType: "XSS",
				Message:     "Possible XSS vulnerability",
				File:        "app/views/users/show.html.erb",
				Line:        10,
				Confidence:  "Medium",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 2)
		require.Equal(t, "SQL Injection", violations[0].CheckName)
		require.Equal(t, "XSS", violations[1].CheckName)
	})

	t.Run("skips invalid warnings and processes valid ones", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
			},
			{
				WarningType: "XSS",
				Message:     "",
				File:        "app/views/users/show.html.erb",
				Line:        10,
				Confidence:  "Medium",
			},
		}

		violations := ConvertWarnings(warnings)
		require.Len(t, violations, 1)
		require.Equal(t, "SQL Injection", violations[0].CheckName)
	})
}

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
