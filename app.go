package awsc

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"github.com/fatih/color"
	"github.com/handlename/awsc/internal/config"
	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

const (
	AWSProfileFlag = "--profile"
	AWSProfileEnv  = "AWS_PROFILE"
)

type App struct {
	config *config.Config
}

func NewApp(configPath string) (*App, error) {
	c, err := config.Load(configPath)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to load config"))
	}

	return &App{
		config: c,
	}, nil
}

func (a *App) Run(ctx context.Context, argv []string) error {
	profile := a.DetectProfile(argv)

	if yes, err := a.ShouldHighlight(profile); err != nil {
		return failure.Wrap(err, failure.Message("failed to check highlight target"))
	} else if yes {
		if err := a.Highlight(profile); err != nil {
			return failure.Wrap(err, failure.Message("failed to hilight"))
		}
	}

	return a.exec(argv)
}

func (a *App) DetectProfile(argv []string) string {
	// read profile from argv
	for i, v := range argv {
		if v == AWSProfileFlag {
			if i+1 < len(argv) {
				p := argv[i+1]
				log.Debug().Str("profile", p).Msg("profile detected from argv")
				return p
			}
		}
	}

	// read profile from envs
	p := os.Getenv(AWSProfileEnv)
	if p != "" {
		log.Debug().Str("profile", p).Msg("profile detected from envs")
		return p
	}

	return ""
}

func (a *App) ShouldHighlight(profile string) (bool, error) {
	if profile == "" {
		log.Debug().Msg("no profile specified. skip to highlight")
		return false, nil
	}

	for _, p := range a.config.Patterns {
		log.Debug().Str("pattern", p).Msg("checking pattern")

		r, err := regexp.Compile(p)
		if err != nil {
			return false, failure.Wrap(err,
				failure.WithCode(errorcode.ErrInvalidArgument),
				failure.Messagef("failed to compile pattern: %s", p))
		}

		if r.MatchString(profile) {
			return true, nil
		}
	}

	log.Debug().Msg("no pattern matched")

	return false, nil
}

func (a *App) Highlight(profile string) error {
	// TODO: ability to change highlight style
	c := color.New(color.FgBlack, color.BgRed)

	if _, err := c.Fprintf(os.Stderr, "%s=%s\n", AWSProfileEnv, profile); err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to highlight"))
	}

	return nil
}

func (a *App) exec(argv []string) error {
	// TODO: ability to change aws cli command path
	cmd := "aws"

	bin, err := exec.LookPath(cmd)
	if err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Messagef("command is not executable %s", cmd))
	}

	args := append([]string{cmd}, argv...)
	log.Debug().Str("bin", bin).Strs("args", args).Msg("exec")

	if err := syscall.Exec(bin, args, os.Environ()); err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to exec aws"))
	}

	return nil
}
