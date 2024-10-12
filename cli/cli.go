package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/handlename/awsc"
	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type ExitCode int

const (
	ExitCodeOK    ExitCode = 0
	ExitCodeError ExitCode = 1
)

func Run() ExitCode {
	logLevel := os.Getenv(awsc.EnvLogLevel)
	if logLevel == "" {
		logLevel = "info"
	}
	awsc.InitLogger(logLevel)

	configPath, err := determineConigPath()
	if err != nil {
		handleError(err)
		return ExitCodeError
	}

	app, err := awsc.NewApp(configPath)
	if err != nil {
		handleError(err)
		return ExitCodeError
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := app.Run(ctx, os.Args[1:]); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Error().Msg("canceled")
		} else {
			handleError(err)
		}

		return ExitCodeError
	}

	return ExitCodeOK
}

func determineConigPath() (string, error) {
	if p := os.Getenv(awsc.EnvConfigPath); p != "" {
		return p, nil
	}

	for _, p := range []string{
		filepath.Join(os.Getenv(awsc.EnvDefaultConfigDir), "awsc", "config.yaml"),
		filepath.Join(os.Getenv("HOME"), ".awsc", "config.yaml"),
	} {
		if _, err := os.Stat(p); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return "", failure.Wrap(err,
				failure.WithCode(errorcode.ErrInternal),
				failure.Message("failed to stat"),
				failure.Context{
					"path": p,
				},
			)
		}

		return p, nil
	}

	return "", nil
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
