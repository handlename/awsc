package entity

import (
	"strings"

	"github.com/handlename/awsc/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type Color int

const (
	Red Color = iota
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Black
	None
)

func ParseColor(color string) (Color, error) {
	if color == "" {
		return Red, nil
	}

	switch strings.ToLower(color) {
	case "red":
		return Red, nil
	case "green":
		return Green, nil
	case "yellow":
		return Yellow, nil
	case "blue":
		return Blue, nil
	case "magenta":
		return Magenta, nil
	case "cyan":
		return Cyan, nil
	case "white":
		return White, nil
	case "black":
		return Black, nil
	default:
		return None, failure.New("invalid color", failure.WithCode(errorcode.ErrInvalidArgument))
	}
}
