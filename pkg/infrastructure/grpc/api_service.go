package grpc

import (
	"connectrpc.com/connect/v2"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
)

var connectCodes = map[domain.ErrorType]connect.Code{
	domain.ErrorTypeBadRequest:         connect.CodeInvalidArgument,
	domain.ErrorTypeNotFound:           connect.CodeNotFound,
	domain.ErrorTypeAlreadyExists:      connect.CodeAlreadyExists,
	domain.ErrorTypeForbidden:          connect.CodePermissionDenied,
	domain.ErrorTypeFailedPrecondition: connect.CodeFailedPrecondition,
}

func handleUseCaseError(err error) error {
	publicMessage, typ, ok := domain.DecomposeError(err)
	if ok {
		if code, known := connectCodes[typ]; known {
			return connect.NewError(code, publicMessage).WithCause(err)
		}
	}
	// Untagged errors are unanticipated: the client gets a generic message, the
	// detail stays in the log.
	return connect.NewError(connect.CodeInternal, "internal server error").WithCause(err)
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
