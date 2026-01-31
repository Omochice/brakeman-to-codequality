package codequality

import (
	"encoding/json"
	"io"
)

type Violation struct {
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

// Write encodes violations as JSON into w.
func Write(violations []Violation, w io.Writer) error {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(violations); err != nil {
		return err
	}
	return nil
}
