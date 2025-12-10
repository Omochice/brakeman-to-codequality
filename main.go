package main

import (
	"fmt"
	"os"
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

func main() {
	fmt.Fprintln(os.Stderr, "brakeman-to-codequality: Not yet implemented")
	os.Exit(1)
}
