package loki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/shiguredo/websocket"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type Config struct {
	Endpoint      string `mapstructure:"endpoint" yaml:"endpoint"`
	QueryTemplate string `mapstructure:"queryTemplate" yaml:"queryTemplate"`
	LogLimit      int    `mapstructure:"logLimit" yaml:"logLimit"`
}

func DefaultQueryTemplate() string {
	return `{ns_trap_jp_app_id="{{ .App.ID }}"}`
}

type lokiStreamer struct {
	config Config
	tmpl   *template.Template
	client *http.Client
}

func NewLokiStreamer(
	config Config,
) (domain.ContainerLogger, error) {
	tmpl, err := template.New("logQL templater").Parse(config.QueryTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "invalid logQL template")
	}

	l := &lokiStreamer{
		config: config,
		tmpl:   tmpl,
		client: &http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost:     30,
				MaxIdleConnsPerHost: 2,
			},
		},
	}

	// check template validity
	var dummy domain.Application
	_, err = l.logQL(&dummy)
	if err != nil {
		return nil, errors.Wrap(err, "executing logQL template")
	}

	return l, nil
}

type m = map[string]any

func templateStr(tmpl *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (l *lokiStreamer) logQL(app *domain.Application) (string, error) {
	return templateStr(l.tmpl, m{"App": app})
}

func (l *lokiStreamer) LogLimit() int {
	return l.config.LogLimit
}

func (l *lokiStreamer) Get(ctx context.Context, app *domain.Application, before time.Time, limit int) ([]*domain.ContainerLog, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", l.queryRangeEndpoint(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	logQL, err := l.logQL(app)
	if err != nil {
		return nil, errors.Wrap(err, "templating logQL")
	}
	q := req.URL.Query()
	q.Set("query", logQL)
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("end", fmt.Sprintf("%d", before.UnixNano()))
	q.Set("since", "1d")
	q.Set("direction", "backward")
	req.URL.RawQuery = q.Encode()

	hres, err := l.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform http request")
	}
	var res queryRangeResponse
	err = json.NewDecoder(hres.Body).Decode(&res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	if res.Status != "success" {
		return nil, errors.Errorf("expected response status to be success, got %s", res.Status)
	}
	if res.Data.ResultType != "streams" {
		return nil, errors.Errorf("expected result type to be streams, got %s", res.Data.ResultType)
	}
	return res.Data.Result.toSortedResponse(true)
}

func (l *lokiStreamer) Stream(ctx context.Context, app *domain.Application, begin time.Time) (<-chan *domain.ContainerLog, error) {
	logQL, err := l.logQL(app)
	if err != nil {
		return nil, errors.Wrap(err, "templating logQL")
	}
	q := make(url.Values)
	q.Set("query", logQL)
	q.Set("limit", "100")
	q.Set("start", fmt.Sprintf("%d", begin.UnixNano()))

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
		defer slog.InfoContext(ctx, "closing loki websocket stream")
		slog.InfoContext(ctx, "new loki websocket stream")

		for {
			typ, b, err := conn.ReadMessage()
			select { // check if context was cancelled
			case <-ctx.Done():
				return
			default:
			}
			if err != nil {
				slog.ErrorContext(ctx, "failed to read ws message", "error", err)
				return
			}
			switch typ {
			case websocket.TextMessage:
				var res streamResponse
				err = json.NewDecoder(bytes.NewReader(b)).Decode(&res)
				if err != nil {
					slog.ErrorContext(ctx, "failed to decode ws message", "error", err)
					continue // fail-safe
				}
				logs, err := res.Streams.toSortedResponse(true)
				if err != nil {
					slog.ErrorContext(ctx, "failed to decode ws message", "error", err)
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
