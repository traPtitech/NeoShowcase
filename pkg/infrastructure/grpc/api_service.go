package grpc

import (
	"github.com/bufbuild/connect-go"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
)

func handleUseCaseError(err error) error {
	underlying, typ, ok := apiserver.DecomposeError(err)
	if ok {
		switch typ {
		case apiserver.ErrorTypeBadRequest:
			return connect.NewError(connect.CodeInvalidArgument, underlying)
		case apiserver.ErrorTypeNotFound:
			return connect.NewError(connect.CodeNotFound, underlying)
		case apiserver.ErrorTypeAlreadyExists:
			return connect.NewError(connect.CodeAlreadyExists, underlying)
		case apiserver.ErrorTypeForbidden:
			return connect.NewError(connect.CodePermissionDenied, underlying)
		}
	}
	return connect.NewError(connect.CodeInternal, err)
}

type APIService struct {
	svc           *apiserver.Service
	avatarBaseURL domain.AvatarBaseURL
}

func NewAPIServiceServer(
	svc *apiserver.Service,
	avatarBaseURL domain.AvatarBaseURL,
) pbconnect.APIServiceHandler {
	return &APIService{
		svc:           svc,
		avatarBaseURL: avatarBaseURL,
	}
}
