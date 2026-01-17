package tarfs

import (
	"archive/tar"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/friendsofgo/errors"
)

func isValidRelPath(relPath string) bool {
	const root = string(filepath.Separator)
	cleaned := filepath.Clean(relPath)
	// filepath.Join cleans traversal path
	// https://dzx.cz/2021-04-02/go_path_traversal/
	traversalCleaned, err := filepath.Rel(root, filepath.Join(root, strings.TrimPrefix(cleaned, root)))
	if err != nil {
		return false
	}
	return traversalCleaned == cleaned
}

func Extract(tarStream io.Reader, destPath string) error {
	tr := tar.NewReader(tarStream)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return errors.Wrap(err, "bad tar file")
		}

		if !isValidRelPath(header.Name) {
			return errors.Errorf("invalid path %v", header.Name)
		}

		path := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(path, header.FileInfo().Mode()); err != nil {
				return errors.Wrap(err, "failed to create directory")
			}

		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Dir(path), header.FileInfo().Mode()|os.ModeDir|100); err != nil {
				return errors.Wrap(err, "failed to create directory")
			}

			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return errors.Wrap(err, "failed to create file")
			}
			_, err = io.Copy(file, tr)
			_ = file.Close()
			if err != nil {
				return errors.Wrap(err, "failed to write file")
			}

		default:
			slog.Debug("skipping tar entry", "name", header.Name, "type", header.Typeflag)
		}
	}
	return nil
}
