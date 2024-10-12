package awsc

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/handlename/awsc/internal/config"
	"github.com/handlename/awsc/internal/entity"
	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

const (
	AWSProfileFlag = "--profile"
	AWSProfileEnv  = "AWS_PROFILE"
)

type App struct {
	config   *config.Config
	patterns []*entity.Pattern
}

func NewApp(configPath string) (*App, error) {
	c, err := config.Load(configPath)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to load config"))
	}

	ps := make([]*entity.Pattern, 0, len(c.Patterns))
	for _, cp := range c.Patterns {
		p, err := entity.NewPattern(cp.Expression, cp.Color)
		if err != nil {
			return nil, failure.Wrap(err, failure.Message("failed to create pattern"))
		}

		ps = append(ps, p)
	}

	return &App{
		config:   c,
		patterns: ps,
	}, nil
}

func (a *App) Run(ctx context.Context, argv []string) error {
	account, err := a.buildAccount(a.detectProfile(argv))
	if err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to build account"))
	}

	if pattern := a.ShouldHighlight(account); pattern != nil {
		if err := a.Highlight(account, pattern); err != nil {
			return failure.Wrap(err, failure.Message("failed to hilight"))
		}
	}

	return a.exec(argv)
}

func (a *App) buildAccount(profile string) (*entity.Account, error) {
	return entity.NewAccount(profile), nil
}

func (a *App) detectProfile(argv []string) string {
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

func (a *App) ShouldHighlight(account *entity.Account) *entity.Pattern {
	if account.Profile() == "" {
		log.Debug().Msg("no profile specified. skip to highlight")
		return nil
	}

	for _, p := range a.patterns {
		if p.Match(account.Profile()) {
			return p
		}
	}

	log.Debug().Msg("no pattern matched")

	return nil
}

var highlightTmpl = template.Must(template.New("highlight").Parse(strings.Join([]string{
	"╓ AWS Account info",
	"╙ Profile: {{ .Profile }}",
}, "\n"),
))

func (a *App) Highlight(account *entity.Account, pattern *entity.Pattern) error {
	var fg color.Attribute

	switch pattern.Color() {
	case entity.Red:
		fg = color.FgRed
	case entity.Green:
		fg = color.FgGreen
	case entity.Yellow:
		fg = color.FgYellow
	case entity.Blue:
		fg = color.FgBlue
	case entity.Magenta:
		fg = color.FgMagenta
	case entity.Cyan:
		fg = color.FgCyan
	case entity.White:
		fg = color.FgWhite
	case entity.Black:
		fg = color.FgBlack
	}

	c := color.New(fg)

	var buf bytes.Buffer
	if err := highlightTmpl.Execute(&buf, map[string]string{
		"Profile": account.Profile(),
	}); err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to execute template"))
	}

	if _, err := c.Fprintf(os.Stderr, buf.String()); err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to highlight"))
	}

	fmt.Println("")

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
