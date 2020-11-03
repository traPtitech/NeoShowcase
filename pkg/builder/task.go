package builder

import (
	"context"
	"github.com/traPtitech/neoshowcase/pkg/models"
)

type Task struct {
	BuildID       string
	RepositoryURL string
	ImageName     string
	Ctx           context.Context
	CancelFunc    func()
	BuildLogM     models.BuildLog
}

func (t *Task) dispose() error {
	if t.CancelFunc != nil {
		t.CancelFunc()
	}
	return nil
}
