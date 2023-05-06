package usecase

import (
	"bytes"
	"sync"

	"github.com/samber/lo"
)

type buildLogStream struct {
	id   string
	b    bytes.Buffer
	subs []chan<- []byte

	lock sync.Mutex
}

type LogStreamService struct {
	buildLogs map[string]*buildLogStream

	lock sync.Mutex
}

func NewLogStreamService() *LogStreamService {
	return &LogStreamService{
		buildLogs: make(map[string]*buildLogStream),
	}
}

func (st *buildLogStream) append(b []byte) {
	st.lock.Lock()
	defer st.lock.Unlock()

	st.b.Write(b)
	for _, sub := range st.subs {
		select {
		case sub <- b:
		default:
		}
	}
}

func (st *buildLogStream) subscribe(sub chan<- []byte) {
	st.lock.Lock()
	defer st.lock.Unlock()

	st.subs = append(st.subs, sub)
	select {
	case sub <- st.b.Bytes():
	default:
	}
}

func (st *buildLogStream) unsubscribe(sub chan<- []byte) {
	st.lock.Lock()
	defer st.lock.Unlock()

	if !lo.Contains(st.subs, sub) {
		return // already closed
	}
	st.subs = lo.Without(st.subs, sub)
	close(sub)
}

func (st *buildLogStream) close() {
	st.lock.Lock()
	defer st.lock.Unlock()

	for _, sub := range st.subs {
		close(sub)
	}
	st.subs = nil
}

func (l *LogStreamService) StartBuildLog(id string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	st, ok := l.buildLogs[id]
	if ok {
		st.close()
	}
	l.buildLogs[id] = &buildLogStream{id: id}
}

func (l *LogStreamService) AppendBuildLog(id string, logPortion []byte) {
	l.lock.Lock()
	st, ok := l.buildLogs[id]
	l.lock.Unlock()

	if ok {
		st.append(logPortion)
	}
}

func (l *LogStreamService) CloseBuildLog(id string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	st, ok := l.buildLogs[id]
	if ok {
		st.close()
	}
	delete(l.buildLogs, id)
}

func (l *LogStreamService) SubscribeBuildLog(id string, sub chan<- []byte) (ok bool, unsubscribe func()) {
	l.lock.Lock()
	st, ok := l.buildLogs[id]
	l.lock.Unlock()

	if !ok {
		return false, nil
	}
	st.subscribe(sub)
	return true, func() {
		st.unsubscribe(sub)
	}
}
