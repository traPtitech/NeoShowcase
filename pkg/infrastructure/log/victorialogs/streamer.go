package victorialogs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo/mutable"

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

type victoriaLogsStreamer struct {
	config Config
	tmpl   *template.Template
	client *http.Client
}

func NewVictoriaLogsStreamer(config Config) (domain.ContainerLogger, error) {
	tmpl, err := template.New("logsQL templater").Parse(config.QueryTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "invalid logsQL template")
	}

	l := &victoriaLogsStreamer{
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
	_, err = l.logsQL(&dummy)
	if err != nil {
		return nil, errors.Wrap(err, "executing logsQL template")
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

func (l *victoriaLogsStreamer) logsQL(app *domain.Application) (string, error) {
	return templateStr(l.tmpl, m{"App": app})
}

func (l *victoriaLogsStreamer) LogLimit() int {
	return l.config.LogLimit
}

func (l *victoriaLogsStreamer) Get(ctx context.Context, app *domain.Application, before time.Time, limit int) ([]*domain.ContainerLog, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", l.queryEndpoint(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	logsQL, err := l.logsQL(app)
	if err != nil {
		return nil, errors.Wrap(err, "templating logsQL")
	}

	start := before.Add(-24 * time.Hour)
	q := req.URL.Query()
	q.Set("query", fmt.Sprintf("_time:[%d, %d) %s | sort by (_time) desc", start.UnixNano(), before.UnixNano(), logsQL))
	q.Set("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()

	res, err := l.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "executing http request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", res.StatusCode, string(body))
	}

	var lines []*domain.ContainerLog
	for line, err := range decodeQuery(res.Body) {
		if err != nil {
			return nil, errors.Wrap(err, "decoding query response")
		}
		lines = append(lines, line)
	}
	mutable.Reverse(lines) // Sort in ascending order
	return lines, nil
}

func (l *victoriaLogsStreamer) Stream(ctx context.Context, app *domain.Application, begin time.Time) (<-chan *domain.ContainerLog, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", l.tailEndpoint(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	logsQL, err := l.logsQL(app)
	if err != nil {
		return nil, errors.Wrap(err, "templating logsQL")
	}

	q := req.URL.Query()
	q.Set("query", fmt.Sprintf("_time:>=%s %s", begin.Format(time.RFC3339), logsQL))
	req.URL.RawQuery = q.Encode()
	res, err := l.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "executing http request")
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", res.StatusCode, string(body))
	}

	ch := make(chan *domain.ContainerLog, 100)

	go func() {
		defer res.Body.Close()
		defer slog.InfoContext(ctx, "closing victorialogs stream")
		slog.InfoContext(ctx, "new victorialogs stream")

		for line, err := range decodeQuery(res.Body) {
			if err != nil {
				slog.ErrorContext(ctx, "failed to decode query response", "error", err)
				return
			}
			select {
			case ch <- line:
			default:
			}
		}
	}()

	return ch, nil
}
