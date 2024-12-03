package builder

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (s *ServiceImpl) saveArtifact(ctx context.Context, st *state) error {
	// Open artifact
	filename := st.artifactTempFile.Name()
	stat, err := os.Stat(filename)
	if err != nil {
		return errors.Wrap(err, "opening artifact")
	}

	// Create artifact meta
	artifact := domain.NewArtifact(st.build.ID, domain.BuilderStaticArtifactName, stat.Size())

	// Create artifact .tar.gz
	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "opening artifact")
	}
	defer file.Close()

	var artifactBytes bytes.Buffer
	gzipWriter := gzip.NewWriter(&artifactBytes)
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		return errors.Wrap(err, "copying file to gzip write")
	}
	err = gzipWriter.Close()
	if err != nil {
		return errors.Wrap(err, "flushing gzip write")
	}

	// Save artifact by requesting to controller
	err = s.client.SaveArtifact(ctx, artifact, artifactBytes.Bytes())
	if err != nil {
		return errors.Wrap(err, "saving artifact")
	}

	return nil
}
