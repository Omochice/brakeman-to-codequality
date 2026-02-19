package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
	"github.com/Omochice/brakeman-to-codequality/cli"
	"github.com/Omochice/brakeman-to-codequality/codequality"
	"github.com/Omochice/brakeman-to-codequality/converter"
)

const version = "0.1.0"

func handleError(w io.Writer, err error) int {
	fmt.Fprintf(w, "Error: %v\n", err)
	return 1
}

func command(args []string, inout *cli.ProcInout) int {
	opts, err := cli.Parse(args)
	if err != nil {
		var helpErr *cli.HelpError
		if errors.As(err, &helpErr) {
			inout.Stderr.Write([]byte(helpErr.Help))
			return 0
		} else {
			return handleError(inout.Stderr, err)
		}
	}

	if opts.Version {
		inout.Stdout.Write([]byte(version))
		return 0
	}

	report, err := brakeman.Parse(inout.Stdin)
	if err != nil {
		return handleError(inout.Stderr, err)
	}

	violations := converter.Warnings(report.Warnings)

	if err := codequality.Write(violations, inout.Stdout); err != nil {
		return handleError(inout.Stderr, err)
	}

	return 0
}

func main() {
	cli.Run(command)
}
