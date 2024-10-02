package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

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
	flags, err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse flags")
		return ExitCodeError
	}

	awsc.InitLogger(flags.LogLevel)

	if flags.Version {
		log.Info().Msgf("awsc v%s", awsc.Version)
		return ExitCodeOK
	}

	app := awsc.NewApp([]string(flags.Patterns), flags.Argv)

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
