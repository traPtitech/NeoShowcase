package loki

import (
	"strconv"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"golang.org/x/exp/slices"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

const (
	queryRangePath = "/loki/api/v1/query_range"
	streamPath     = "/loki/api/v1/tail"
)

func (l *lokiStreamer) queryRangeEndpoint() string {
	return l.config.Endpoint + queryRangePath
}

func (l *lokiStreamer) streamEndpoint() string {
	wsEndpoint := l.config.Endpoint
	wsEndpoint = strings.ReplaceAll(wsEndpoint, "http://", "ws://")
	wsEndpoint = strings.ReplaceAll(wsEndpoint, "https://", "wss://")
	return wsEndpoint + streamPath
}

type queryRangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string       `json:"resultType"`
		Result     streamValues `json:"result"`
		// Stats map[string]any `json:"stats"`
	} `json:"data"`
}

type streamValue struct {
	Stream map[string]string `json:"stream"`
	Values []value           `json:"values"`
}

type value [2]string

func (v *value) time() (time.Time, error) {
	unixNano, err := strconv.Atoi(v[0])
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, int64(unixNano)), nil
}

func (v *value) log() string {
	return v[1]
}

type streamResponse struct {
	Streams streamValues `json:"streams"`
	// DroppedEntries []droppedEntry `json:"dropped_entries"`
}

type streamValues []streamValue

func (values streamValues) toSortedResponse(asc bool) ([]*domain.ContainerLog, error) {
	var logs []*domain.ContainerLog
	for _, sv := range values {
		for _, v := range sv.Values {
			logTime, err := v.time()
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode response time")
			}
			logs = append(logs, &domain.ContainerLog{
				Time: logTime,
				Log:  v.log(),
			})
		}
	}
	if asc {
		slices.SortFunc(logs, ds.LessFunc(func(a *domain.ContainerLog) int64 { return a.Time.UnixNano() }))
	} else {
		slices.SortFunc(logs, ds.MoreFunc(func(a *domain.ContainerLog) int64 { return a.Time.UnixNano() }))
	}
	return logs, nil
}
