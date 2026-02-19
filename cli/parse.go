package cli

import (
	"bytes"

	"github.com/jessevdk/go-flags"
)

func Parse(args []string) (*Options, error) {
	var opts Options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := parser.ParseArgs(args)
	if err != nil {
		if ferr, ok := err.(*flags.Error); ok && ferr.Type == flags.ErrHelp {
			var buf bytes.Buffer
			parser.WriteHelp(&buf)
			return nil, NewHelpError(buf.String())
		}
		return nil, err
	}
	return &opts, err
}
