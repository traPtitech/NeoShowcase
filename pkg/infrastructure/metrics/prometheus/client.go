package prometheus

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const defaultStep = 1 * time.Minute

type QueryConfig struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Template string `mapstructure:"template" yaml:"template"`
}

func DefaultQueriesConfig() []*QueryConfig {
	// default values for k8s cadvisor
	selector := `namespace="ns-apps", pod="nsapp-{{ .App.ID }}-0", container="app"`
	return []*QueryConfig{
		{
			Name:     "CPU",
			Template: fmt.Sprintf(`rate(container_cpu_user_seconds_total{%s}[5m]) + rate(container_cpu_system_seconds_total{%s}[5m])`, selector, selector),
		},
		{
			Name:     "Memory",
			Template: fmt.Sprintf(`container_memory_usage_bytes{%s} + container_memory_swap{%s}`, selector, selector),
		},
	}
}

type Config struct {
	Endpoint string         `mapstructure:"endpoint" yaml:"endpoint"`
	Queries  []*QueryConfig `mapstructure:"queries" yaml:"queries"`
}

type promClient struct {
	config    Config
	templates map[string]*template.Template
	client    promv1.API
}

func NewPromClient(
	config Config,
) (domain.MetricsService, error) {
	templates := make(map[string]*template.Template, len(config.Queries))
	for _, qc := range config.Queries {
		tmpl, err := template.New(fmt.Sprintf("promQL templater %v", qc.Name)).Parse(qc.Template)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("invalid promQL template %v", qc.Name))
		}
		templates[qc.Name] = tmpl
	}

	client, err := api.NewClient(api.Config{Address: config.Endpoint})
	if err != nil {
		return nil, errors.Wrap(err, "creating prom cleint")
	}
	p := &promClient{
		config:    config,
		templates: templates,
		client:    promv1.NewAPI(client),
	}

	// check templates validity
	var dummy domain.Application
	for _, qc := range config.Queries {
		_, err = p.promQL(qc.Name, &dummy)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("executing logQL template %v", qc.Name))
		}
	}

	return p, nil
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

func (p *promClient) promQL(name string, app *domain.Application) (string, error) {
	tmpl, ok := p.templates[name]
	if !ok {
		return "", errors.Errorf("no such template: %v", name)
	}
	return templateStr(tmpl, m{"App": app})
}

func (p *promClient) AvailableNames() []string {
	return ds.Map(p.config.Queries, func(qc *QueryConfig) string { return qc.Name })
}

func (p *promClient) Get(ctx context.Context, name string, app *domain.Application, before time.Time, limit time.Duration) ([]*domain.AppMetric, error) {
	promQL, err := p.promQL(name, app)
	if err != nil {
		return nil, errors.Wrap(err, "templating promQL")
	}
	v, _, err := p.client.QueryRange(ctx, promQL, promv1.Range{
		Start: before.Add(-limit),
		End:   before.Add(1), // to inclusive
		Step:  defaultStep,
	})
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}

	if v.Type() != model.ValMatrix {
		return nil, errors.Errorf("expected result type to be matrix, but got %v", v.Type().String())
	}
	mv, ok := v.(model.Matrix)
	if !ok {
		return nil, errors.New("cast value failed")
	}
	return toSortedResponse(mv), nil
}

func toSortedResponse(mv model.Matrix) []*domain.AppMetric {
	var items []*domain.AppMetric
	for _, item := range mv {
		for _, v := range item.Values {
			items = append(items, &domain.AppMetric{
				Time:  v.Timestamp.Time(),
				Value: float64(v.Value),
			})
		}
	}
	slices.SortFunc(items, ds.LessFunc(func(e *domain.AppMetric) int64 { return e.Time.UnixNano() }))
	return items
}
