package entity

import (
	"regexp"

	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type Rule struct {
	expression      *regexp.Regexp
	color           Color
	confirmOnModify bool
}

func NewRule(expression, color string, confirm bool) (*Rule, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return &Rule{}, failure.Wrap(err,
			failure.WithCode(errorcode.ErrInvalidArgument),
			failure.Message("failed to compile expression"),
			failure.Context{
				"expression": expression,
				"color":      color,
			})
	}

	c, err := ParseColor(color)
	if err != nil {
		return &Rule{}, failure.Wrap(err,
			failure.WithCode(errorcode.ErrInvalidArgument),
			failure.Message("failed to parse color"),
			failure.Context{
				"expression": expression,
				"color":      color,
			})
	}

	return &Rule{
		expression:      r,
		color:           c,
		confirmOnModify: confirm,
	}, nil
}

func (p *Rule) Color() Color {
	return p.color
}

func (p *Rule) Match(s string) bool {
	return p.expression.MatchString(s)
}

func (p *Rule) ConfirmOnModify() bool {
	return p.confirmOnModify
}
