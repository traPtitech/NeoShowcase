package isolatedrunner

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
)

// ReceiveStepName is the name of the step that each implementation runs before the steps of the job to extract its input.
const ReceiveStepName = "receive"

var (
	jobIDPattern    = regexp.MustCompile(`^[a-z0-9]{1,32}$`)
	stepNamePattern = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$`)
)

// ValidateJob reports whether the job satisfies the preconditions documented on builder.IsolatedJob.
func ValidateJob(job builder.IsolatedJob) error {
	if !jobIDPattern.MatchString(job.ID) {
		return oops.Errorf("invalid job id %q", job.ID)
	}
	if job.Image == "" {
		return oops.New("job image is empty")
	}
	if job.Input == nil {
		return oops.New("job input is nil")
	}
	if len(job.Steps) == 0 {
		return oops.New("job has no steps")
	}
	for _, dir := range job.Shared {
		if !strings.HasPrefix(dir, "/") {
			return oops.Errorf("shared directory %q is not absolute", dir)
		}
	}
	names := map[string]bool{ReceiveStepName: true}
	for _, step := range job.Steps {
		if !stepNamePattern.MatchString(step.Name) || names[step.Name] {
			return oops.Errorf("invalid or duplicate step name %q", step.Name)
		}
		names[step.Name] = true
		if len(step.Cmd) == 0 {
			return oops.Errorf("step %s has no command", step.Name)
		}
	}
	return nil
}

// ReceiveScript returns a shell script, run as root, that extracts the input tar from stdin at /
// and makes the shared directories owned by RunAs.
func ReceiveScript(job builder.IsolatedJob) string {
	script := "tar x -C /"
	if len(job.Shared) == 0 {
		return script
	}
	quoted := make([]string, len(job.Shared))
	for i, dir := range job.Shared {
		quoted[i] = shellQuote(dir)
	}
	return script + fmt.Sprintf(" && chown -R %d:%d %s", job.RunAs.UID, job.RunAs.GID, strings.Join(quoted, " "))
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// WriteStepStart writes the line marking the start of a step to the build log.
func WriteStepStart(logWriter io.Writer, stepName string) {
	_, _ = fmt.Fprintf(logWriter, "[ns-builder] Running step %q\n", stepName)
}
