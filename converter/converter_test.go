package converter_test

import (
	"testing"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
	"github.com/Omochice/brakeman-to-codequality/converter"
)

func TestSeverity(t *testing.T) {
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
			got := converter.Severity(tt.confidence)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWarnings(t *testing.T) {
	t.Run("converts valid warning correctly", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
				Code:        "User.where(...)",
				Fingerprint: "a21418b38aa77ef73946105fb1c9e3623b7be67a2419b960793871587200cbcc",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(violations))
		}

		violation := violations[0]
		if violation.Description != "Possible SQL injection" {
			t.Fatalf("got %v, want %v", violation.Description, "Possible SQL injection")
		}
		if violation.CheckName != "SQL Injection" {
			t.Fatalf("got %v, want %v", violation.CheckName, "SQL Injection")
		}
		if violation.Severity != "critical" {
			t.Fatalf("got %v, want %v", violation.Severity, "critical")
		}
		if violation.Location.Path != "app/models/user.rb" {
			t.Fatalf("got %v, want %v", violation.Location.Path, "app/models/user.rb")
		}
		if violation.Location.Lines.Begin != 42 {
			t.Fatalf("got %v, want %v", violation.Location.Lines.Begin, 42)
		}
		if violation.Fingerprint != "a21418b38aa77ef73946105fb1c9e3623b7be67a2419b960793871587200cbcc" {
			t.Fatalf("got %v, want %v", violation.Fingerprint, "a21418b38aa77ef73946105fb1c9e3623b7be67a2419b960793871587200cbcc")
		}
	})

	t.Run("skips warning with missing file", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "",
				Line:        42,
				Confidence:  "High",
				Fingerprint: "abc123",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips warning with missing line", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        0,
				Confidence:  "High",
				Fingerprint: "abc123",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips warning with missing warning type", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
				Fingerprint: "abc123",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips warning with missing message", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
				Fingerprint: "abc123",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips warning with missing fingerprint", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("removes ./ prefix from file path", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "./app/models/user.rb",
				Line:        42,
				Confidence:  "High",
				Fingerprint: "abc123",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(violations))
		}
		if violations[0].Location.Path != "app/models/user.rb" {
			t.Fatalf("got %v, want %v", violations[0].Location.Path, "app/models/user.rb")
		}
	})

	t.Run("handles empty array", func(t *testing.T) {
		warnings := []brakeman.Warning{}
		violations := converter.Warnings(warnings)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("processes multiple warnings", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
				Fingerprint: "fp1",
			},
			{
				WarningType: "XSS",
				Message:     "Possible XSS vulnerability",
				File:        "app/views/users/show.html.erb",
				Line:        10,
				Confidence:  "Medium",
				Fingerprint: "fp2",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 2 {
			t.Fatalf("expected length %d, got %d", 2, len(violations))
		}
		if violations[0].CheckName != "SQL Injection" {
			t.Fatalf("got %v, want %v", violations[0].CheckName, "SQL Injection")
		}
		if violations[1].CheckName != "XSS" {
			t.Fatalf("got %v, want %v", violations[1].CheckName, "XSS")
		}
	})

	t.Run("skips invalid warnings and processes valid ones", func(t *testing.T) {
		warnings := []brakeman.Warning{
			{
				WarningType: "SQL Injection",
				Message:     "Possible SQL injection",
				File:        "app/models/user.rb",
				Line:        42,
				Confidence:  "High",
				Fingerprint: "fp1",
			},
			{
				WarningType: "XSS",
				Message:     "",
				File:        "app/views/users/show.html.erb",
				Line:        10,
				Confidence:  "Medium",
				Fingerprint: "fp2",
			},
		}

		violations := converter.Warnings(warnings)
		if len(violations) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(violations))
		}
		if violations[0].CheckName != "SQL Injection" {
			t.Fatalf("got %v, want %v", violations[0].CheckName, "SQL Injection")
		}
	})
}
