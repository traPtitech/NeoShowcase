package fig

import (
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/lukesampson/figlet/figletlib"
)

type Builder struct {
	lines []string
}

func (b *Builder) Append(msg string, font string, color Color) error {
	by, err := readFontBytes(font)
	if err != nil {
		return nil
	}
	f, err := figletlib.ReadFontFromBytes(by)
	if err != nil {
		return errors.Wrap(err, "reading font from bytes")
	}
	words := figletlib.GetLines(msg, f, 1000, f.Settings())
	for _, word := range words {
		art := word.Art()
		if len(b.lines) > 0 && len(b.lines) != len(art) {
			return errors.New("line height mismatch")
		}
		if len(b.lines) == 0 {
			b.lines = make([]string, len(art))
		}

		for i, line := range art {
			if color != nil {
				b.lines[i] += color.GetPrefix()
			}
			b.lines[i] += string(line)
			if color != nil {
				b.lines[i] += color.GetSuffix()
			}
		}
	}
	return nil
}

func (b *Builder) String() string {
	return strings.Join(b.lines, "\n")
}
