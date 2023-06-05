package grpc

import (
	"github.com/bufbuild/connect-go"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
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
	svc           *usecase.APIServerService
	avatarBaseURL domain.AvatarBaseURL
}

func NewAPIServiceServer(
	svc *usecase.APIServerService,
	avatarBaseURL domain.AvatarBaseURL,
) pbconnect.APIServiceHandler {
	return &APIService{
		svc:           svc,
		avatarBaseURL: avatarBaseURL,
	}
}
