package storage

import (
	"fmt"
	"io"
	"os"

	"github.com/ncw/swift"
)

type SwiftStorage struct {
	Container string
	Conn      *swift.Connection
}

func (ss *SwiftStorage) Save(filename string, src io.Reader) error {
	_, err := ss.Conn.ObjectPut(ss.Container, filename, src, true, "", "", swift.Headers{})
	return err
}

func (ss *SwiftStorage) Open(filename string) (io.ReadCloser, error) {
	file, _, err := ss.Conn.ObjectOpen(ss.Container, filename, true, nil)
	if err != nil {
		if err == swift.ObjectNotFound {
			return nil, fmt.Errorf("not found: %w", err)
		}
		return nil, err
	}
	return file, nil
}

func (ss *SwiftStorage) Delete(filename string) error {
	err := ss.Conn.ObjectDelete(ss.Container, filename)
	if err != nil {
		if err == swift.ObjectNotFound {
			return fmt.Errorf("not found: %w", err)
		}
		return err
	}
	return nil
}

func (ss *SwiftStorage) Move(sourcePath, destPath string) error {
	// Move LocalDir to Swift Storage
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	_, err = ss.Conn.ObjectPut(ss.Container, destPath, inputFile, true, "", "", swift.Headers{})
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %w", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %w", err)
	}
	return nil
}

func (ss *SwiftStorage) FileExists(filename string) bool {
	_, _, err := ss.Conn.Object(ss.Container, filename)
	return err == nil
}
