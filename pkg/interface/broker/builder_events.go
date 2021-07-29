package broker

import (
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BuilderEventsBroker interface {
	Run() error
}

type builderEventsBroker struct {
	bus    domain.Bus
	client pb.BuilderServiceClient
}

func NewBuilderEventsBroker(client pb.BuilderServiceClient, bus domain.Bus) (BuilderEventsBroker, error) {
	return &builderEventsBroker{
		bus:    bus,
		client: client,
	}, nil
}

func (b *builderEventsBroker) Run() error {
	stream, err := b.client.ConnectEventStream(context.Background(), &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to connect events tream: %w", err)
	}

	for {
		ev, err := stream.Recv()
		if err == io.EOF {
			log.Debug("builder event stream was closed: EOF")
			break
		}
		if err != nil {
			log.WithError(err).
				Debug("builder event stream was disconnected with error")
			return err
		}

		payload := util.FromJSON(ev.Body)

		log.WithField("type", ev.Type).
			WithField("payload", payload).
			Info("builder event received")

		switch ev.Type {
		case pb.Event_BUILD_STARTED:
			b.bus.Publish(event.BuilderBuildStarted, payload)
		case pb.Event_BUILD_SUCCEEDED:
			b.bus.Publish(event.BuilderBuildSucceeded, payload)
		case pb.Event_BUILD_FAILED:
			b.bus.Publish(event.BuilderBuildFailed, payload)
		case pb.Event_BUILD_CANCELED:
			b.bus.Publish(event.BuilderBuildCanceled, payload)
		}
	}
	return nil
}
