package grpc

import (
	"context"
	"errors"
	"io"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/samber/lo"
	"github.com/samber/oops"

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
		client: pbconnect.NewBuildpackHelperServiceClient(connect.NewClient(connecthttp.NewTransport(web.NewH2CClient(), address, connecthttp.WithReadMaxBytes(0)))),
	}
}

func (b *BuildpackHelperServiceClient) CopyFileTree(ctx context.Context, destination string, tarStream io.Reader) error {
	content, err := io.ReadAll(tarStream)
	if err != nil {
		return err
	}
	req := &pb.CopyFileTreeRequest{
		Destination: destination,
		TarContent:  content,
	}
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
	req := &pb.HelperExecRequest{
		WorkDir: workDir,
		Cmd:     cmd,
		Envs:    arrayEnvs,
	}
	st, err := b.client.Exec(ctx, req)
	if err != nil {
		return 0, oops.Wrapf(err, "requesting exec")
	}
	var exitCode *int
	for {
		msg, err := st.Receive()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return 0, oops.Wrapf(err, "receiving logs")
		}
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

	if exitCode == nil {
		return 0, oops.New("exit code not received")
	}
	return *exitCode, nil
}
