package brakeman_test

import (
	"strings"
	"testing"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("parses valid JSON", func(t *testing.T) {
		input := `{"warnings":[{"warning_type":"SQL Injection","message":"Possible SQL injection","file":"app/models/user.rb","line":42,"confidence":"High"}]}`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		require.NoError(t, err)
		require.NotNil(t, report)
		require.Len(t, report.Warnings, 1)
		require.Equal(t, "SQL Injection", report.Warnings[0].WarningType)
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		input := `{invalid json`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		require.Error(t, err)
		require.Nil(t, report)
	})

	t.Run("handles empty warnings array", func(t *testing.T) {
		input := `{"warnings":[]}`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		require.NoError(t, err)
		require.NotNil(t, report)
		require.Len(t, report.Warnings, 0)
	})

	t.Run("handles missing warnings field", func(t *testing.T) {
		input := `{}`
		reader := strings.NewReader(input)

		report, err := brakeman.Parse(reader)
		require.NoError(t, err)
		require.NotNil(t, report)
		require.Len(t, report.Warnings, 0)
	})
}
