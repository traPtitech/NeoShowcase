package fig

// Ported from https://github.com/mbndr/figlet4go/blob/master/color.go

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// Escape char
const escape string = "\x1b"

// Terminal AnsiColors
var (
	ColorBlack   = AnsiColor{30}
	ColorRed     = AnsiColor{31}
	ColorGreen   = AnsiColor{32}
	ColorYellow  = AnsiColor{33}
	ColorBlue    = AnsiColor{34}
	ColorMagenta = AnsiColor{35}
	ColorCyan    = AnsiColor{36}
	ColorWhite   = AnsiColor{37}
)

// TrueColorForAnsiColor is TrueColor lookalikes for displaying AnsiColor f.e. with the HTML parser
// Colors based on http://clrs.cc/
var TrueColorForAnsiColor = map[AnsiColor]TrueColor{
	ColorBlack:   {0, 0, 0},
	ColorRed:     {255, 65, 54},
	ColorGreen:   {149, 189, 64},
	ColorYellow:  {255, 220, 0},
	ColorBlue:    {0, 116, 217},
	ColorMagenta: {177, 13, 201},
	ColorCyan:    {105, 206, 245},
	ColorWhite:   {255, 255, 255},
}

// Color has a pre- and a suffix
type Color interface {
	// GetPrefix returns prefix for ansi color
	GetPrefix() string
	// GetSuffix returns suffix for ansi color
	GetSuffix() string
}

// AnsiColor representation
type AnsiColor struct {
	code int
}

// TrueColor with rgb Attributes
type TrueColor struct {
	r int
	g int
	b int
}

func (tc TrueColor) GetPrefix() string {
	return fmt.Sprintf("%v[38;2;%d;%d;%dm", escape, tc.r, tc.g, tc.b)
}

func (tc TrueColor) GetSuffix() string {
	return fmt.Sprintf("%v[0m", escape)
}

// NewTrueColorFromHexString returns a TrueColor object based on a hexadecimal string
func NewTrueColorFromHexString(c string) (*TrueColor, error) {
	rgb, err := hex.DecodeString(c)
	if err != nil {
		return nil, errors.New("Invalid color given (" + c + ")")
	}

	return &TrueColor{
		int(rgb[0]),
		int(rgb[1]),
		int(rgb[2]),
	}, nil
}

func (ac AnsiColor) GetPrefix() string {
	return fmt.Sprintf("%v[0;%dm", escape, ac.code)
}

func (ac AnsiColor) GetSuffix() string {
	return fmt.Sprintf("%v[0m", escape)
}
