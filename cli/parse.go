package cli

import (
	"github.com/jessevdk/go-flags"
)

func Parse(args *[]string) (*Options, []string, error) {
	var opts Options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	a, err := parser.ParseArgs(*args)
	if err != nil {
		return nil, nil, err
	}
	return &opts, a, err
}
