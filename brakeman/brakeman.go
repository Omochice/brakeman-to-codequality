package brakeman

import (
	"encoding/json"
	"io"
)

type Report struct {
	Warnings []Warning `json:"warnings"`
}

type Warning struct {
	WarningType string `json:"warning_type"`
	Message     string `json:"message"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Confidence  string `json:"confidence"`
	Code        string `json:"code,omitempty"`
}

// Parse decodes a Brakeman JSON report from r.
func Parse(r io.Reader) (*Report, error) {
	var report Report

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&report); err != nil {
		return nil, err
	}

	if report.Warnings == nil {
		report.Warnings = []Warning{}
	}

	return &report, nil
}
