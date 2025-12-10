package main

import (
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
