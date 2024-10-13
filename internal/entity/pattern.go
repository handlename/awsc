package entity

import (
	"regexp"

	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type Rule struct {
	expression *regexp.Regexp
	color      Color
}

func NewRule(expression, color string) (*Rule, error) {
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

	return &Rule{expression: r, color: c}, nil
}

func (p *Rule) Color() Color {
	return p.color
}

func (p *Rule) Match(s string) bool {
	return p.expression.MatchString(s)
}
