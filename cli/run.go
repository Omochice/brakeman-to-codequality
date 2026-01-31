package cli

import (
	"io"
	"os"
)

// ProcInout abstracts standard I/O streams so that the main pipeline
// can be tested without capturing os.Stdin/os.Stdout/os.Stderr.
type ProcInout struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewProcInout() *ProcInout {
	return &ProcInout{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

type Command func(args []string, inout *ProcInout) int

// Run executes a Command with the real process arguments and I/O, then exits.
func Run(c Command) {
	args := os.Args[1:]
	exitStatus := c(args, NewProcInout())
	os.Exit(exitStatus)
}
