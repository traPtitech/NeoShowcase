package fig

import (
	"embed"
	"io"
	"path/filepath"

	"github.com/friendsofgo/errors"
)

//go:embed fonts/*
var fs embed.FS

const fontsDir = "fonts"

func readFontBytes(font string) ([]byte, error) {
	f, err := fs.Open(filepath.Join(fontsDir, font+".flf"))
	if err != nil {
		return nil, errors.Wrap(err, "opening font file")
	}
	defer f.Close()
	return io.ReadAll(f)
}
