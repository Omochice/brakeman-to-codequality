package cli

import "testing"

func TestParse(t *testing.T) {
	t.Run("returns error when no positional argument given", func(t *testing.T) {
		_, err := Parse([]string{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error when multiple positional arguments given", func(t *testing.T) {
		_, err := Parse([]string{"a.json", "b.json"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("sets Source to the positional argument", func(t *testing.T) {
		opts, err := Parse([]string{"report.json"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.Source != "report.json" {
			t.Fatalf("got %q, want %q", opts.Source, "report.json")
		}
	})

	t.Run("sets Source to dash for stdin", func(t *testing.T) {
		opts, err := Parse([]string{"-"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.Source != "-" {
			t.Fatalf("got %q, want %q", opts.Source, "-")
		}
	})

	t.Run("allows no positional argument when version flag is set", func(t *testing.T) {
		opts, err := Parse([]string{"--version"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !opts.Version {
			t.Fatal("expected Version to be true")
		}
	})
}
