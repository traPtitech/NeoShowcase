package tarfs

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
)

func Compress(srcPath string) io.Reader {
	pr, pw := io.Pipe()
	tw := tar.NewWriter(pw)

	go func() {
		defer pw.Close()

		err := filepath.Walk(srcPath, func(file string, fi os.FileInfo, err error) error {
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(srcPath, file)
			if err != nil {
				return err
			}
			if rel == "." {
				return nil // skip root dir
			}
			header.Name = rel
			if fi.IsDir() {
				header.Name += "/"
			}

			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}
			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}
				_, err = io.Copy(tw, data)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			_ = pw.CloseWithError(errors.Wrap(err, "walking srcPath"))
		}
		if err = tw.Close(); err != nil {
			_ = pw.CloseWithError(errors.Wrap(err, "closing tar writer"))
		}
	}()

	return pr
}
