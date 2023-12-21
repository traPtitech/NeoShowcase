package grpc

import (
	"connectrpc.com/connect"
	"context"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"io"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
)

type BuildpackHelperServiceClient struct {
	client pbconnect.BuildpackHelperServiceClient
}

func NewBuildpackHelperServiceClient(
	address string,
) domain.BuildpackHelperServiceClient {
	return &BuildpackHelperServiceClient{
		client: pbconnect.NewBuildpackHelperServiceClient(web.NewH2CClient(), address),
	}
}

func (b *BuildpackHelperServiceClient) CopyFileTree(ctx context.Context, destination string, tarStream io.Reader) error {
	content, err := io.ReadAll(tarStream)
	if err != nil {
		return err
	}
	req := connect.NewRequest(&pb.CopyFileTreeRequest{
		Destination: destination,
		TarContent:  content,
	})
	_, err = b.client.CopyFileTree(ctx, req)
	return err
}

func (b *BuildpackHelperServiceClient) Exec(ctx context.Context, workDir string, cmd []string, envs map[string]string, logWriter io.Writer) (int, error) {
	arrayEnvs := lo.MapToSlice(envs, func(key string, value string) *pb.HelperExecEnv {
		return &pb.HelperExecEnv{
			Key:   key,
			Value: value,
		}
	})
	req := connect.NewRequest(&pb.HelperExecRequest{
		WorkDir: workDir,
		Cmd:     cmd,
		Envs:    arrayEnvs,
	})
	st, err := b.client.Exec(ctx, req)
	if err != nil {
		return 0, errors.Wrap(err, "requesting exec")
	}
	var exitCode *int
	for st.Receive() {
		msg := st.Msg()
		switch msg.Type {
		case pb.HelperExecResponse_LOG:
			payload := msg.Body.(*pb.HelperExecResponse_Log).Log
			_, err = logWriter.Write(payload)
			if err != nil {
				return 0, err
			}
		case pb.HelperExecResponse_EXIT_CODE:
			payload := int(msg.Body.(*pb.HelperExecResponse_ExitCode).ExitCode)
			exitCode = &payload
		}
	}
	if err := st.Err(); err != nil {
		return 0, errors.Wrap(err, "receiving logs")
	}
	if exitCode == nil {
		return 0, errors.New("exit code not received")
	}
	return *exitCode, nil
}
