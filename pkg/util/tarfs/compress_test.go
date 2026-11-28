package tarfs

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddDir_WritesSymlinksAsLinks(t *testing.T) {
	src := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(src, "src"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(src, "src", "main.go"), []byte("package main"), 0o644))
	require.NoError(t, os.Symlink("src/main.go", filepath.Join(src, "link")))

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	require.NoError(t, AddDir(tw, src, "app", nil))
	require.NoError(t, tw.Close())

	entries := map[string]*tar.Header{}
	contents := map[string]string{}
	tr := tar.NewReader(&buf)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		b, err := io.ReadAll(tr)
		require.NoError(t, err)
		entries[h.Name] = h
		contents[h.Name] = string(b)
	}

	require.Contains(t, entries, "app/")
	require.Equal(t, "package main", contents["app/src/main.go"])
	link := entries["app/link"]
	require.NotNil(t, link)
	require.Equal(t, byte(tar.TypeSymlink), link.Typeflag)
	require.Equal(t, "src/main.go", link.Linkname)
	require.Empty(t, contents["app/link"])
}
