package cli

import (
	"flag"
	"fmt"

	"github.com/morikuni/failure/v2"
)

type Flags struct {
	Version  bool
	LogLevel string
	Patterns Patterns
	Argv     []string
}

type Patterns []string

func (p *Patterns) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *Patterns) Set(pattern string) error {
	*p = append(*p, pattern)
	return nil
}

func parseFlags(appname string, argv []string) (*Flags, error) {
	flags := &Flags{}

	fs := flag.NewFlagSet(appname, flag.ExitOnError)

	fs.BoolVar(&flags.Version, "version", false, "Print version")
	fs.StringVar(&flags.LogLevel, "log-level", "info", "Log level (trace, debug, info, warn, error, panic)")
	fs.Var(&flags.Patterns, "pattern", "Pattern to match")

	if err := fs.Parse(argv); err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to parse flags"))
	}

	flags.Argv = fs.Args()

	return flags, nil
}
