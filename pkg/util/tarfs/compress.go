package tarfs

import (
	"archive/tar"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/samber/oops"
)

// AddDir writes the tree under srcPath to tw, naming each entry by its path relative to srcPath joined to prefix.
// srcPath itself is written only when prefix is not empty. Symbolic links are written as links.
// If edit is not nil, each header is passed to it before being written, and may be modified there.
func AddDir(tw *tar.Writer, srcPath string, prefix string, edit func(*tar.Header)) error {
	err := filepath.Walk(srcPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var link string
		if fi.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(file)
			if err != nil {
				return err
			}
		}
		header, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcPath, file)
		if err != nil {
			return err
		}
		if rel == "." && prefix == "" {
			return nil // skip root dir
		}
		header.Name = path.Join(prefix, filepath.ToSlash(rel))
		if fi.IsDir() {
			header.Name += "/"
		}
		if edit != nil {
			edit(header)
		}

		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()
			_, err = io.Copy(tw, data)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return oops.Wrapf(err, "walking srcPath")
	}
	return nil
}
