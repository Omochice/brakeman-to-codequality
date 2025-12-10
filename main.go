package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// BrakemanReport represents the structure of Brakeman JSON output
type BrakemanReport struct {
	Warnings []BrakemanWarning `json:"warnings"`
}

// BrakemanWarning represents a single security warning from Brakeman
type BrakemanWarning struct {
	WarningType string `json:"warning_type"`
	Message     string `json:"message"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Confidence  string `json:"confidence"`
	Code        string `json:"code,omitempty"` // Optional field
}

// CodeQualityViolation represents a GitLab Code Quality violation
type CodeQualityViolation struct {
	Description string   `json:"description"`
	CheckName   string   `json:"check_name"`
	Fingerprint string   `json:"fingerprint"`
	Severity    string   `json:"severity"`
	Location    Location `json:"location"`
}

// Location represents the file location of a violation
type Location struct {
	Path  string `json:"path"`
	Lines Lines  `json:"lines"`
}

// Lines represents the line numbers where a violation occurs
type Lines struct {
	Begin int `json:"begin"`
}

// MapSeverity converts Brakeman confidence levels to GitLab severity values
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

// GenerateFingerprint creates a unique SHA-256 hash for a warning
func GenerateFingerprint(file string, line int, warningType, message, code string) string {
	// Combine fields to create input for hashing
	input := file + ":" + strconv.Itoa(line) + ":" + warningType + ":" + message
	if code != "" {
		input += ":" + code
	}

	// Calculate SHA-256 hash
	hash := sha256.Sum256([]byte(input))

	// Convert to 64-character hex string
	return hex.EncodeToString(hash[:])
}

// ConvertWarnings transforms Brakeman warnings to GitLab Code Quality violations
func ConvertWarnings(warnings []BrakemanWarning) []CodeQualityViolation {
	violations := make([]CodeQualityViolation, 0, len(warnings))

	for _, warning := range warnings {
		// Skip warnings with missing required fields
		if warning.File == "" || warning.Line == 0 || warning.WarningType == "" || warning.Message == "" {
			continue
		}

		// Remove "./" prefix from file path
		path := strings.TrimPrefix(warning.File, "./")

		// Create violation
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

// ParseBrakemanJSON reads and parses Brakeman JSON from an io.Reader
func ParseBrakemanJSON(r io.Reader) (*BrakemanReport, error) {
	var report BrakemanReport

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&report); err != nil {
		return nil, err
	}

	// If warnings field is nil, initialize as empty slice
	if report.Warnings == nil {
		report.Warnings = []BrakemanWarning{}
	}

	return &report, nil
}

func main() {
	fmt.Fprintln(os.Stderr, "brakeman-to-codequality: Not yet implemented")
	os.Exit(1)
}
