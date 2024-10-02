package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/handlename/awsc"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type ExitCode int

const (
	ExitCodeOK    ExitCode = 0
	ExitCodeError ExitCode = 1
)

func Run() ExitCode {
	var cli struct {
		Version  bool     `help:"Print version"`
		LogLevel string   `help:"Log level" enum:"trace,debug,info,warn,error,panic" env:"AWSC_LOG_LEVEL" default:"info"`
		Patterns []string `help:"Pattern for AWS profile to highlight" name:"pattern" env:"AWSC_PATTERN" default:"production"`
		Argv     []string `arg:""`
	}

	kc := kong.Parse(&cli)
	if err := kc.Error; err != nil {
		log.Error().Err(err).Msg("failed to parse flags")
		return ExitCodeError
	}

	awsc.InitLogger(cli.LogLevel)

	if cli.Version {
		log.Info().Msgf("awsc v%s", awsc.Version)
		return ExitCodeOK
	}

	log.Debug().Strs("args", cli.Argv)
	app := awsc.NewApp(cli.Patterns, cli.Argv)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := app.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Error().Msg("canceled")
		} else {
			handleError(err)
		}

		return ExitCodeError
	}

	return ExitCodeOK
}

func handleError(err error) {
	fmt.Println("======== error ========")

	code := failure.CodeOf(err)
	fmt.Printf("code = %s\n", code)

	msg := failure.MessageOf(err)
	fmt.Printf("message = %s\n", msg)

	cs := failure.CallStackOf(err)
	fmt.Printf("callstack = %s\n", cs)

	fmt.Printf("cause = %s\n", failure.CauseOf(err))

	fmt.Println()
	fmt.Println("======== detail ========")
	fmt.Printf("%+v\n", err)
}
