package builder

import (
	"compress/gzip"
	"context"
	"io"
	"os"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (s *builderService) saveArtifact(ctx context.Context, st *state) error {
	filename := st.artifactTempFile.Name()

	stat, err := os.Stat(filename)
	if err != nil {
		return errors.Wrap(err, "opening artifact")
	}

	artifact := domain.NewArtifact(st.build.ID, domain.BuilderStaticArtifactName, stat.Size())
	err = s.artifactRepo.CreateArtifact(ctx, artifact)
	if err != nil {
		return errors.Wrap(err, "creating artifact record")
	}

	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "opening artifact")
	}
	defer file.Close()

	pr, pw := io.Pipe()
	gzipWriter := gzip.NewWriter(pw)
	if err != nil {
		return errors.Wrap(err, "creating gzip stream")
	}
	go func() {
		defer pw.Close()
		_, err := io.Copy(gzipWriter, file)
		if err != nil {
			_ = pw.CloseWithError(errors.Wrap(err, "copying file to pipe writer"))
			return
		}
		err = gzipWriter.Close()
		if err != nil {
			_ = pw.CloseWithError(errors.Wrap(err, "flushing gzip writer"))
			return
		}
	}()
	err = domain.SaveArtifact(s.storage, artifact.ID, pr)
	if err != nil {
		return errors.Wrap(err, "saving artifact")
	}

	return nil
}
