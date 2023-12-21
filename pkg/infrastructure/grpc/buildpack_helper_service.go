package grpc

import (
	"bytes"
	"context"
	"github.com/friendsofgo/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"sync"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/tarfs"
)

type BuildpackHelperService struct {
	ssGenConnections []*ssGenConnection
	lock             sync.Mutex
}

func NewBuildpackHelperService() pbconnect.BuildpackHelperServiceHandler {
	return &BuildpackHelperService{}
}

func (b *BuildpackHelperService) CopyFileTree(_ context.Context, req *connect.Request[pb.CopyFileTreeRequest]) (*connect.Response[emptypb.Empty], error) {
	err := tarfs.Extract(bytes.NewReader(req.Msg.TarContent), req.Msg.Destination)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&emptypb.Empty{})
	return res, nil
}

type buildpackHelperExec struct {
	st *connect.ServerStream[pb.HelperExecResponse]
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

func (b *BuildpackHelperService) Exec(ctx context.Context, req *connect.Request[pb.HelperExecRequest], st *connect.ServerStream[pb.HelperExecResponse]) error {
	if len(req.Msg.Cmd) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("cmd cannot have length of 0"))
	}

	// Prepare command
	cmd := exec.CommandContext(ctx, req.Msg.Cmd[0], req.Msg.Cmd[1:]...)
	cmd.Dir = req.Msg.WorkDir
	cmd.Env = os.Environ()
	additionalEnvs := ds.Map(req.Msg.Envs, func(env *pb.HelperExecEnv) string { return env.Key + "=" + env.Value })
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
	var exitError *exec.ExitError
	if errors.As(cmdErr, &exitError) {
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
