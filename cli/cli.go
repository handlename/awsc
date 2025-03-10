package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/handlename/awsc"
	"github.com/handlename/awsc/internal/env"
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
	if s := os.Getenv(env.EnvShowVersion); s != "" {
		fmt.Printf("awsc v%s", awsc.Version)
		return ExitCodeOK
	}

	logLevel := os.Getenv(env.EnvLogLevel)
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
	if p := os.Getenv(env.EnvConfigPath); p != "" {
		log.Debug().Str("path", p).Msg("using config path from environment variable")
		return p, nil
	}

	for _, p := range []string{
		filepath.Join(os.Getenv(env.EnvDefaultConfigDir), "awsc", "config.yaml"),
		filepath.Join(os.Getenv("HOME"), ".config", "awsc", "config.yaml"),
		filepath.Join(os.Getenv("HOME"), ".awsc", "config.yaml"),
	} {
		log.Debug().Str("path", p).Msg("checking config path")
		if _, err := os.Stat(p); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Debug().Str("path", p).Msg("config path not found")
				continue
			}

			var pathErr *fs.PathError
			if errors.As(err, &pathErr) {
				// When a file is placed in path than expects a directory
				if strings.Contains(pathErr.Error(), "not a directory") {
					return "", failure.Wrap(err,
						failure.WithCode(errorcode.ErrInternal),
						failure.Message("a file is placed in path than expects a directory"),
						failure.Context{
							"path": p,
						},
					)
				}
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
