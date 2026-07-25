package grpc

import (
	"connectrpc.com/connect"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
)

// publicError reports only the client-facing message from Error(), which is what
// Connect puts on the wire, while still unwrapping to the internal chain so the
// boundary log keeps the full detail.
type publicError struct {
	public string
	cause  error
}

func (e publicError) Error() string { return e.public }
func (e publicError) Unwrap() error { return e.cause }

var connectCodes = map[apiserver.ErrorType]connect.Code{
	apiserver.ErrorTypeBadRequest:    connect.CodeInvalidArgument,
	apiserver.ErrorTypeNotFound:      connect.CodeNotFound,
	apiserver.ErrorTypeAlreadyExists: connect.CodeAlreadyExists,
	apiserver.ErrorTypeForbidden:     connect.CodePermissionDenied,
}

func handleUseCaseError(err error) error {
	publicMessage, typ, ok := apiserver.DecomposeError(err)
	if ok {
		if code, known := connectCodes[typ]; known {
			return connect.NewError(code, publicError{public: publicMessage, cause: err})
		}
	}
	// Untagged errors are unanticipated: the client gets a generic message, the
	// detail stays in the log.
	return connect.NewError(connect.CodeInternal, publicError{
		public: "internal server error",
		cause:  err,
	})
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
