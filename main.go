package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
	"github.com/Omochice/brakeman-to-codequality/cli"
	"github.com/Omochice/brakeman-to-codequality/codequality"
)

func MapSeverity(confidence string) string {
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

func GenerateFingerprint(file string, line int, warningType, message, code string) string {
	input := file + ":" + strconv.Itoa(line) + ":" + warningType + ":" + message
	if code != "" {
		input += ":" + code
	}

	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func ConvertWarnings(warnings []brakeman.Warning) []codequality.Violation {
	violations := make([]codequality.Violation, 0, len(warnings))

	for _, warning := range warnings {
		if warning.File == "" || warning.Line == 0 || warning.WarningType == "" || warning.Message == "" {
			continue
		}

		path := strings.TrimPrefix(warning.File, "./")

		violation := codequality.Violation{
			Description: warning.Message,
			CheckName:   warning.WarningType,
			Fingerprint: GenerateFingerprint(warning.File, warning.Line, warning.WarningType, warning.Message, warning.Code),
			Severity:    MapSeverity(warning.Confidence),
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

func handleError(w io.Writer, err error) int {
	fmt.Fprintf(w, "Error: %v\n", err)
	return 1
}

func command(args []string, inout *cli.ProcInout) int {
	report, err := brakeman.Parse(inout.Stdin)
	if err != nil {
		return handleError(inout.Stderr, err)
	}

	violations := ConvertWarnings(report.Warnings)

	if err := codequality.Write(violations, inout.Stdout); err != nil {
		return handleError(inout.Stderr, err)
	}

	return 0
}

func main() {
	cli.Run(command)
}
