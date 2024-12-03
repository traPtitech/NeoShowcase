package pbconnect

import (
	"context"
	"errors"
	http "net/http"
	"testing"

	connect "connectrpc.com/connect"
	// "github.com/traPtitech/neoshowcase/pkg/domain"
	pb "github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
)

type aaa struct {
	a string
}

func (aaa) Do(*http.Request) (*http.Response, error) {
	return &(http.Response{}), nil
}

func setUp() (aPIServiceClient, error) {
	// setting original gateway and builder for testing
	connect1 := aaa{}
	// a := connect.HTTPClient{}
	ac := NewAPIServiceClient(connect1, "")

	asc, ok := ac.(*aPIServiceClient)
	if !ok {
		return aPIServiceClient{}, errors.New("おしまいです")
	}

	return *asc, nil
}

func TestBuildRuntimeCmd(t *testing.T) {
	gw, err := setUp()
	if err != nil {
		return
	}

	gw.RetryCommitBuild(context.TODO(), &connect.Request[pb.RetryCommitBuildRequest]{})

	// throw appId used to identify

	// client.RegisterBuild(context.TODO(), "")

	// not table driven
	// tests := []struct {
	// 	name    string
	// 	sr      domain.StartBuildRequest
	// 	wantErr bool
	// }{}

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		// err := tt.bs.startBuild(&tt.sr)
	// 	})
	// }

	// 最後に err チェックをするテストを書く〜
	// assert.NoError(t, err)
}
