package usecase

import (
	"bytes"
	"sync"

	"github.com/samber/lo"
)

type LogStreamService struct {
	buildLogs map[string]*buildLogStream

	lock sync.Mutex
}

func NewLogStreamService() *LogStreamService {
	return &LogStreamService{
		buildLogs: make(map[string]*buildLogStream),
	}
}

type buildLogStream struct {
	id   string
	b    bytes.Buffer
	subs []chan<- []byte
}

func (st *buildLogStream) append(b []byte) {
	st.b.Write(b)
	for _, sub := range st.subs {
		select {
		case sub <- b:
		default:
		}
	}
}

func (st *buildLogStream) subscribe(sub chan<- []byte) {
	st.subs = append(st.subs, sub)
	select {
	case sub <- st.b.Bytes():
	default:
	}
}

func (st *buildLogStream) unsubscribe(sub chan<- []byte) {
	st.subs = lo.Without(st.subs, sub)
	close(sub)
}

func (st *buildLogStream) close() {
	for _, sub := range st.subs {
		close(sub)
	}
}

func (l *LogStreamService) getOrInitBuildLogStream(id string) *buildLogStream {
	st, ok := l.buildLogs[id]
	if !ok {
		st = &buildLogStream{id: id}
		l.buildLogs[id] = st
	}
	return st
}

func (l *LogStreamService) AppendBuildLog(id string, logPortion []byte) {
	l.lock.Lock()
	defer l.lock.Unlock()

	st := l.getOrInitBuildLogStream(id)
	st.append(logPortion)
}

func (l *LogStreamService) CloseBuildLog(id string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	st, ok := l.buildLogs[id]
	if ok {
		st.close()
	}
}

func (l *LogStreamService) SubscribeBuildLog(id string, sub chan<- []byte) (ok bool, unsubscribe func()) {
	l.lock.Lock()
	defer l.lock.Unlock()

	st, ok := l.buildLogs[id]
	if !ok {
		return false, nil
	}
	st.subscribe(sub)
	return true, func() {
		st.unsubscribe(sub)
	}
}
