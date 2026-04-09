package brakeman_test

import (
	"strings"
	"testing"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
)

func TestParse(t *testing.T) {
	t.Run("parses valid JSON", func(t *testing.T) {
		input := `{"warnings":[{"warning_type":"SQL Injection","message":"Possible SQL injection","file":"app/models/user.rb","line":42,"confidence":"High"}]}`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if report == nil {
			t.Fatalf("expected non-nil value")
		}
		if len(report.Warnings) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(report.Warnings))
		}
		if report.Warnings[0].WarningType != "SQL Injection" {
			t.Fatalf("got %v, want %v", report.Warnings[0].WarningType, "SQL Injection")
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		input := `{invalid json`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if report != nil {
			t.Fatalf("expected nil, got %v", report)
		}
	})

	t.Run("handles empty warnings array", func(t *testing.T) {
		input := `{"warnings":[]}`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if report == nil {
			t.Fatalf("expected non-nil value")
		}
		if len(report.Warnings) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(report.Warnings))
		}
	})

	t.Run("handles missing warnings field", func(t *testing.T) {
		input := `{}`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if report == nil {
			t.Fatalf("expected non-nil value")
		}
		if len(report.Warnings) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(report.Warnings))
		}
	})
}
