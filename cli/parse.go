package cli

import (
	"bytes"
	"fmt"

	"github.com/jessevdk/go-flags"
)

func Parse(args []string) (*Options, error) {
	var opts Options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	parser.Usage = "[OPTIONS] <file path|->"
	remaining, err := parser.ParseArgs(args)
	if err != nil {
		if ferr, ok := err.(*flags.Error); ok && ferr.Type == flags.ErrHelp {
			var buf bytes.Buffer
			parser.WriteHelp(&buf)
			return nil, NewHelpError(buf.String())
		}
		return nil, err
	}

	if opts.Version {
		return &opts, nil
	}

	if len(remaining) != 1 {
		return nil, fmt.Errorf("exactly one argument required (file path or \"-\" for stdin), got %d", len(remaining))
	}
	opts.Source = remaining[0]

	return &opts, nil
}
