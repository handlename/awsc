package entity

import (
	"regexp"

	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type Pattern struct {
	expression *regexp.Regexp
	color      Color
}

func NewPattern(expression, color string) (*Pattern, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return &Pattern{}, failure.Wrap(err,
			failure.WithCode(errorcode.ErrInvalidArgument),
			failure.Message("failed to compile expression"),
			failure.Context{
				"expression": expression,
				"color":      color,
			})
	}

	c, err := ParseColor(color)
	if err != nil {
		return &Pattern{}, failure.Wrap(err,
			failure.WithCode(errorcode.ErrInvalidArgument),
			failure.Message("failed to parse color"),
			failure.Context{
				"expression": expression,
				"color":      color,
			})
	}

	return &Pattern{expression: r, color: c}, nil
}

func (p *Pattern) Color() Color {
	return p.color
}

func (p *Pattern) Match(s string) bool {
	return p.expression.MatchString(s)
}
