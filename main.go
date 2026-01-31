package main

import (
	"fmt"
	"io"

	"github.com/Omochice/brakeman-to-codequality/brakeman"
	"github.com/Omochice/brakeman-to-codequality/cli"
	"github.com/Omochice/brakeman-to-codequality/codequality"
	"github.com/Omochice/brakeman-to-codequality/converter"
)

func handleError(w io.Writer, err error) int {
	fmt.Fprintf(w, "Error: %v\n", err)
	return 1
}

func command(args []string, inout *cli.ProcInout) int {
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
