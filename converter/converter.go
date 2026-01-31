package converter

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
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

// Fingerprint generates a deterministic SHA-256 hex digest for a warning.
func Fingerprint(file string, line int, warningType, message, code string) string {
	input := file + ":" + strconv.Itoa(line) + ":" + warningType + ":" + message
	if code != "" {
		input += ":" + code
	}

	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Warnings converts Brakeman warnings into CodeQuality violations.
// Warnings that lack a file, line, warning type, or message are skipped.
func Warnings(warnings []brakeman.Warning) []codequality.Violation {
	violations := make([]codequality.Violation, 0, len(warnings))

	for _, warning := range warnings {
		if warning.File == "" || warning.Line == 0 || warning.WarningType == "" || warning.Message == "" {
			continue
		}

		path := strings.TrimPrefix(warning.File, "./")

		violation := codequality.Violation{
			Description: warning.Message,
			CheckName:   warning.WarningType,
			Fingerprint: Fingerprint(warning.File, warning.Line, warning.WarningType, warning.Message, warning.Code),
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
