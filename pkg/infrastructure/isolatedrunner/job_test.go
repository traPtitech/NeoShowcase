package isolatedrunner

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

func TestValidateJob(t *testing.T) {
	valid := func() builder.IsolatedJob {
		return builder.IsolatedJob{
			ID:     "abc123",
			Image:  "alpine:3.22",
			Shared: []string{"/data"},
			Input:  strings.NewReader(""),
			Steps:  []builder.IsolatedStep{{Name: "build", Cmd: []string{"true"}}, {Name: "export", Cmd: []string{"true"}}},
		}
	}
	tests := []struct {
		name    string
		modify  func(j *builder.IsolatedJob)
		wantErr string
	}{
		{name: "valid", modify: func(*builder.IsolatedJob) {}},
		{name: "uppercase id", modify: func(j *builder.IsolatedJob) { j.ID = "ABC123" }, wantErr: "job id"},
		{name: "id longer than 32", modify: func(j *builder.IsolatedJob) { j.ID = strings.Repeat("a", 33) }, wantErr: "job id"},
		{name: "no image", modify: func(j *builder.IsolatedJob) { j.Image = "" }, wantErr: "image"},
		{name: "relative shared directory", modify: func(j *builder.IsolatedJob) { j.Shared = []string{"data"} }, wantErr: "not absolute"},
		{name: "no input", modify: func(j *builder.IsolatedJob) { j.Input = nil }, wantErr: "input"},
		{name: "no steps", modify: func(j *builder.IsolatedJob) { j.Steps = nil }, wantErr: "no steps"},
		{name: "step named receive", modify: func(j *builder.IsolatedJob) { j.Steps[0].Name = ReceiveStepName }, wantErr: "step name"},
		{name: "duplicate step names", modify: func(j *builder.IsolatedJob) { j.Steps[1].Name = "build" }, wantErr: "step name"},
		{name: "step name not a DNS label", modify: func(j *builder.IsolatedJob) { j.Steps[0].Name = "Build_1" }, wantErr: "step name"},
		{name: "empty command", modify: func(j *builder.IsolatedJob) { j.Steps[0].Cmd = nil }, wantErr: "no command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := valid()
			tt.modify(&job)
			err := ValidateJob(job)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}
