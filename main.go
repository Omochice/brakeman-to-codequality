package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Omochice/brakeman-to-codequality/cli"
)

type BrakemanReport struct {
	Warnings []BrakemanWarning `json:"warnings"`
}

type BrakemanWarning struct {
	WarningType string `json:"warning_type"`
	Message     string `json:"message"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Confidence  string `json:"confidence"`
	Code        string `json:"code,omitempty"`
}

type CodeQualityViolation struct {
	Description string   `json:"description"`
	CheckName   string   `json:"check_name"`
	Fingerprint string   `json:"fingerprint"`
	Severity    string   `json:"severity"`
	Location    Location `json:"location"`
}

type Location struct {
	Path  string `json:"path"`
	Lines Lines  `json:"lines"`
}

type Lines struct {
	Begin int `json:"begin"`
}

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

func ConvertWarnings(warnings []BrakemanWarning) []CodeQualityViolation {
	violations := make([]CodeQualityViolation, 0, len(warnings))

	for _, warning := range warnings {
		if warning.File == "" || warning.Line == 0 || warning.WarningType == "" || warning.Message == "" {
			continue
		}

		path := strings.TrimPrefix(warning.File, "./")

		violation := CodeQualityViolation{
			Description: warning.Message,
			CheckName:   warning.WarningType,
			Fingerprint: GenerateFingerprint(warning.File, warning.Line, warning.WarningType, warning.Message, warning.Code),
			Severity:    MapSeverity(warning.Confidence),
			Location: Location{
				Path: path,
				Lines: Lines{
					Begin: warning.Line,
				},
			},
		}

		violations = append(violations, violation)
	}

	return violations
}

func ParseBrakemanJSON(r io.Reader) (*BrakemanReport, error) {
	var report BrakemanReport

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&report); err != nil {
		return nil, err
	}

	if report.Warnings == nil {
		report.Warnings = []BrakemanWarning{}
	}

	return &report, nil
}

func WriteCodeQualityJSON(violations []CodeQualityViolation, w io.Writer) error {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(violations); err != nil {
		return err
	}
	return nil
}

func handleError(w io.Writer, err error) int {
	fmt.Fprintf(w, "Error: %v\n", err)
	return 1
}

func command(args []string, inout *cli.ProcInout) int {
	report, err := ParseBrakemanJSON(inout.Stdin)
	if err != nil {
		return handleError(inout.Stderr, err)
	}

	violations := ConvertWarnings(report.Warnings)

	if err := WriteCodeQualityJSON(violations, inout.Stdout); err != nil {
		return handleError(inout.Stderr, err)
	}

	return 0
}

func main() {
	cli.Run(command)
}
