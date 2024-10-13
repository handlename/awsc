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
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
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
	config *config.Config
	rules  []*entity.Rule
}

func NewApp(configPath string) (*App, error) {
	log.Debug().Str("config", configPath).Msg("load config")
	c, err := config.Load(configPath)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to load config"))
	}

	ps := make([]*entity.Rule, 0, len(c.Rules))
	for _, cp := range c.Rules {
		p, err := entity.NewRule(cp.Expression, cp.Color)
		if err != nil {
			return nil, failure.Wrap(err, failure.Message("failed to create rule"))
		}

		ps = append(ps, p)
	}

	return &App{
		config: c,
		rules:  ps,
	}, nil
}

func (a *App) Run(ctx context.Context, argv []string) error {
	account, err := a.buildAccount(ctx, a.detectProfile(argv))
	if err != nil {
		return failure.Wrap(err,
			failure.WithCode(errorcode.ErrInternal),
			failure.Message("failed to build account"))
	}

	if rule := a.ShouldHighlight(account); rule != nil {
		if err := a.Highlight(account, rule); err != nil {
			return failure.Wrap(err, failure.Message("failed to hilight"))
		}
	}

	return a.exec(argv)
}

func (a *App) buildAccount(ctx context.Context, profile string) (*entity.Account, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to load aws config"))
	}

	opts := []entity.AccountOption{}
	if a.config.AdditionalInfo {
		opts = append(opts, entity.AccountOptionWithAdditionalInfo)
	}

	account, err := entity.NewAccount(ctx, profile, cfg, opts...)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to create account"))
	}

	return account, nil
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

func (a *App) ShouldHighlight(account *entity.Account) *entity.Rule {
	if account.Profile() == "" {
		log.Debug().Msg("no profile specified. skip to highlight")
		return nil
	}

	for _, p := range a.rules {
		if p.Match(account.Profile()) {
			return p
		}
	}

	log.Debug().Msg("no rule matched")

	return nil
}

var defaultTmpl = template.Must(template.New("default").Parse(strings.Join([]string{
	"╓ AWS Account info",
	"║ Profile {{ .Profile }}",
	"║ Region  {{ .Region }}",
	`║ ID      {{ if .ID }}{{ .ID }}{{ else }}N/A{{ end }}`,
	`║ ARN     {{ if .Arn }}{{ .Arn }}{{ else }}N/A{{ end }}`,
	`║ UserID  {{ if .UserID }}{{ .UserID }}{{ else }}N/A{{ end }}`,
	"╙ [{{ .Now }}]",
}, "\n"),
))

func (a *App) Highlight(account *entity.Account, rule *entity.Rule) error {
	var fg color.Attribute

	switch rule.Color() {
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

	tmpl := defaultTmpl
	if a.config.Template != "" {
		var err error
		tmpl, err = template.New("custom").Parse(a.config.Template)
		if err != nil {
			return failure.Wrap(err,
				failure.WithCode(errorcode.ErrInvalidArgument),
				failure.Message("failed to parse template"))
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"Now": time.Now().Format(a.config.TimeFormat),

		"Profile": account.Profile(),
		"Region":  account.Region(),

		// Additional info
		"ID":     account.ID(),
		"UserID": account.UserID(),
		"Arn":    account.Arn(),
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
