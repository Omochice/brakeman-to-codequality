package converter

import (
	"strings"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
	"github.com/Omochice/brakeman-to-codequality/codequality"
)

// Severity maps a Brakeman confidence level to a CodeQuality severity.
func Severity(confidence string) string {
	switch strings.ToLower(confidence) {
	case "high":
		return "critical"
	case "medium":
		return "major"
	case "weak", "low":
		return "minor"
	default:
		return "info"
	}
}

// Warnings converts Brakeman warnings into CodeQuality violations.
// Warnings that lack a file, line, warning type, message, or fingerprint are skipped.
func Warnings(warnings []brakeman.Warning) []codequality.Violation {
	violations := make([]codequality.Violation, 0, len(warnings))

	for _, warning := range warnings {
		if warning.File == "" || warning.Line == 0 || warning.WarningType == "" || warning.Message == "" || warning.Fingerprint == "" {
			continue
		}

		path := strings.TrimPrefix(warning.File, "./")

		violation := codequality.Violation{
			Description: warning.Message,
			CheckName:   warning.WarningType,
			Fingerprint: warning.Fingerprint,
			Severity:    Severity(warning.Confidence),
			Location: codequality.Location{
				Path: path,
				Lines: codequality.Lines{
					Begin: warning.Line,
				},
			},
		}

		violations = append(violations, violation)
	}

	return violations
}
