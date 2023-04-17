package grpc

import (
	"github.com/bufbuild/connect-go"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func handleUseCaseError(err error) error {
	underlying, typ, ok := usecase.DecomposeError(err)
	if ok {
		switch typ {
		case usecase.ErrorTypeBadRequest:
			return connect.NewError(connect.CodeInvalidArgument, underlying)
		case usecase.ErrorTypeNotFound:
			return connect.NewError(connect.CodeNotFound, underlying)
		case usecase.ErrorTypeAlreadyExists:
			return connect.NewError(connect.CodeAlreadyExists, underlying)
		case usecase.ErrorTypeForbidden:
			return connect.NewError(connect.CodePermissionDenied, underlying)
		}
	}
	return connect.NewError(connect.CodeInternal, err)
}

type APIService struct {
	svc    *usecase.APIServerService
	pubKey *ssh.PublicKeys
}

func NewAPIServiceServer(
	svc *usecase.APIServerService,
	pubKey *ssh.PublicKeys,
) pbconnect.APIServiceHandler {
	return &APIService{
		svc:    svc,
		pubKey: pubKey,
	}
}
