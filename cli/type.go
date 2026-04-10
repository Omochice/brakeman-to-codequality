package cli

type Options struct {
	Version bool   `short:"v" long:"version" description:"Show application version"`
	Source  string // positional argument: file path or "-" for stdin
}
