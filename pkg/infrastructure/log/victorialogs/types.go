package victorialogs

import (
	"encoding/json"
	"io"
	"iter"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

const (
	queryPath = "/select/logsql/query"
	tailPath  = "/select/logsql/tail"
)

func (l *victoriaLogsStreamer) queryEndpoint() string {
	return l.config.Endpoint + queryPath
}

func (l *victoriaLogsStreamer) tailEndpoint() string {
	return l.config.Endpoint + tailPath
}

type logLine struct {
	Time time.Time `json:"_time"`
	Msg  string    `json:"_msg"`
}

func decodeQuery(r io.Reader) iter.Seq2[*domain.ContainerLog, error] {
	return func(yield func(*domain.ContainerLog, error) bool) {
		decoder := json.NewDecoder(r)
		for decoder.More() {
			var line logLine
			if err := decoder.Decode(&line); err != nil {
				if err == io.EOF {
					break
				}
				yield(nil, err)
				return
			}
			domainLog := &domain.ContainerLog{
				Time: line.Time,
				Log:  line.Msg,
			}
			if !yield(domainLog, nil) {
				return
			}
		}
	}
}
