package main

import (
	"regexp"
	"testing"

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
