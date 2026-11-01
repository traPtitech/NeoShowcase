package grpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"

	"connectrpc.com/connect/v2"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/tarfs"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BuildpackHelperService struct{}

func NewBuildpackHelperService() pbconnect.BuildpackHelperServiceHandler {
	return &BuildpackHelperService{}
}

func (b *BuildpackHelperService) CopyFileTree(_ context.Context, req *pb.CopyFileTreeRequest) (*emptypb.Empty, error) {
	err := tarfs.Extract(bytes.NewReader(req.TarContent), req.Destination)
	if err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return res, nil
}

type buildpackHelperExec struct {
	st pbconnect.BuildpackHelperServiceExecServerStream
}

var _ io.Writer = (*buildpackHelperExec)(nil)

func (b *buildpackHelperExec) Write(p []byte) (n int, err error) {
	err = b.st.Send(&pb.HelperExecResponse{
		Type: pb.HelperExecResponse_LOG,
		Body: &pb.HelperExecResponse_Log{Log: p},
	})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (b *BuildpackHelperService) Exec(ctx context.Context, req *pb.HelperExecRequest, st pbconnect.BuildpackHelperServiceExecServerStream) error {
	if len(req.Cmd) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, "cmd cannot have length of 0")
	}

	// Prepare command
	cmd := exec.CommandContext(ctx, req.Cmd[0], req.Cmd[1:]...)
	cmd.Dir = req.WorkDir
	cmd.Env = os.Environ()
	additionalEnvs := ds.Map(req.Envs, func(env *pb.HelperExecEnv) string { return env.Key + "=" + env.Value })
	cmd.Env = append(cmd.Env, additionalEnvs...) // Inherit important envs such as CNB_STACK_ID

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	// Run command and send logs
	var eg errgroup.Group
	eg.Go(func() error {
		dst := &buildpackHelperExec{st: st}
		_, err := io.Copy(dst, pr)
		return err
	})
	cmdErr := cmd.Run()
	_ = pw.Close()
	logErr := eg.Wait()
	if logErr != nil {
		return logErr
	}

	// Check exit code
	if exitError, ok := errors.AsType[*exec.ExitError](cmdErr); ok {
		return st.Send(&pb.HelperExecResponse{
			Type: pb.HelperExecResponse_EXIT_CODE,
			Body: &pb.HelperExecResponse_ExitCode{ExitCode: int32(exitError.ExitCode())},
		})
	}
	if cmdErr != nil {
		return cmdErr
	}
	return st.Send(&pb.HelperExecResponse{
		Type: pb.HelperExecResponse_EXIT_CODE,
		Body: &pb.HelperExecResponse_ExitCode{ExitCode: 0},
	})
}
