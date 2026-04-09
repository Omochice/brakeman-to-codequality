package converter_test

import (
	"regexp"
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

func TestFingerprint(t *testing.T) {
	t.Run("generates consistent fingerprint for same input", func(t *testing.T) {
		fp1 := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "User.where(...)")
		fp2 := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "User.where(...)")
		if fp1 != fp2 {
			t.Fatalf("got %v, want %v", fp1, fp2)
		}
	})

	t.Run("generates different fingerprints for different inputs", func(t *testing.T) {
		fp1 := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := converter.Fingerprint("app/models/user.rb", 43, "SQL Injection", "Possible SQL injection", "")
		if fp1 == fp2 {
			t.Fatalf("expected values to differ, but both are %v", fp1)
		}
	})

	t.Run("includes code field when present", func(t *testing.T) {
		fpWithCode := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "User.where(...)")
		fpWithoutCode := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		if fpWithCode == fpWithoutCode {
			t.Fatalf("expected values to differ, but both are %v", fpWithCode)
		}
	})

	t.Run("handles empty code field", func(t *testing.T) {
		fp := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		if fp == "" {
			t.Fatalf("expected non-empty string")
		}
	})

	t.Run("generates 64-character hex string", func(t *testing.T) {
		fp := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		if len(fp) != 64 {
			t.Fatalf("expected length %d, got %d", 64, len(fp))
		}
		re := regexp.MustCompile("^[0-9a-f]{64}$")
		if !re.MatchString(fp) {
			t.Fatalf("expected %q to match pattern %v", fp, re)
		}
	})

	t.Run("different files produce different fingerprints", func(t *testing.T) {
		fp1 := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := converter.Fingerprint("app/models/post.rb", 42, "SQL Injection", "Possible SQL injection", "")
		if fp1 == fp2 {
			t.Fatalf("expected values to differ, but both are %v", fp1)
		}
	})

	t.Run("different warning types produce different fingerprints", func(t *testing.T) {
		fp1 := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := converter.Fingerprint("app/models/user.rb", 42, "XSS", "Possible SQL injection", "")
		if fp1 == fp2 {
			t.Fatalf("expected values to differ, but both are %v", fp1)
		}
	})

	t.Run("different messages produce different fingerprints", func(t *testing.T) {
		fp1 := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Possible SQL injection", "")
		fp2 := converter.Fingerprint("app/models/user.rb", 42, "SQL Injection", "Confirmed SQL injection", "")
		if fp1 == fp2 {
			t.Fatalf("expected values to differ, but both are %v", fp1)
		}
	})
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
		if violation.Fingerprint == "" {
			t.Fatalf("expected non-empty string")
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
			},
			{
				WarningType: "XSS",
				Message:     "Possible XSS vulnerability",
				File:        "app/views/users/show.html.erb",
				Line:        10,
				Confidence:  "Medium",
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
			},
			{
				WarningType: "XSS",
				Message:     "",
				File:        "app/views/users/show.html.erb",
				Line:        10,
				Confidence:  "Medium",
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
