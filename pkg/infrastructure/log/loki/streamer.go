package loki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/shiguredo/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type Config struct {
	Endpoint   string `mapstructure:"endpoint" yaml:"endpoint"`
	AppIDLabel string `mapstructure:"appIDLabel" yaml:"appIDLabel"`
}

type lokiStreamer struct {
	config Config
	client *http.Client
}

func NewLokiStreamer(
	config Config,
) domain.ContainerLogger {
	return &lokiStreamer{
		config: config,
		client: &http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost:     30,
				MaxIdleConnsPerHost: 2,
			},
		},
	}
}

func (l *lokiStreamer) logQL(appID string) string {
	return fmt.Sprintf("{%s=\"%s\"}", l.config.AppIDLabel, appID)
}

func (l *lokiStreamer) Stream(ctx context.Context, appID string) (<-chan *domain.ContainerLog, error) {
	q := make(url.Values)
	q.Set("query", l.logQL(appID))
	q.Set("limit", "100")
	q.Set("start", fmt.Sprintf("%d", time.Now().Add(-startLimit).UnixNano()))

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, l.streamEndpoint()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial to stream ws endpoint")
	}

	ch := make(chan *domain.ContainerLog, 100)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-ctx.Done()
		_ = conn.Close()
		defer close(ch)
	}()
	go func() {
		defer cancel()
		defer log.Infof("closing loki websocket stream")
		log.Infof("new loki websocket stream")

		for {
			typ, b, err := conn.ReadMessage()
			select { // check if context was cancelled
			case <-ctx.Done():
				return
			default:
			}
			if err != nil {
				log.Errorf("failed to read ws message: %+v", err)
				return
			}
			switch typ {
			case websocket.TextMessage:
				var res streamResponse
				err = json.NewDecoder(bytes.NewReader(b)).Decode(&res)
				if err != nil {
					log.Errorf("failed to decode ws message: %+v", err)
					continue // fail-safe
				}
				logs, err := res.Streams.toSortedResponse(true)
				if err != nil {
					log.Errorf("failed to decode ws message: %+v", err)
					continue // fail-safe
				}
				for _, l := range logs {
					select {
					case ch <- l:
					default:
					}
				}
			case websocket.BinaryMessage:
				// ignore
			case websocket.CloseMessage:
				return
			}
		}
	}()

	return ch, nil
}
