package logstream

import (
	"sync"
)

type Service struct {
	buildLogs map[string]*buildLogStream

	lock sync.Mutex
}

func NewService() *Service {
	return &Service{
		buildLogs: make(map[string]*buildLogStream),
	}
}

func (l *Service) StartBuildLog(id string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	st, ok := l.buildLogs[id]
	if ok {
		st.close()
	}
	l.buildLogs[id] = &buildLogStream{id: id}
}

func (l *Service) AppendBuildLog(id string, logPortion []byte) {
	l.lock.Lock()
	st, ok := l.buildLogs[id]
	l.lock.Unlock()

	if ok {
		st.append(logPortion)
	}
}

func (l *Service) CloseBuildLog(id string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	st, ok := l.buildLogs[id]
	if ok {
		st.close()
	}
	delete(l.buildLogs, id)
}

func (l *Service) HasBuildLog(id string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	_, ok := l.buildLogs[id]
	return ok
}

func (l *Service) SubscribeBuildLog(id string, sub chan<- []byte) (ok bool, unsubscribe func()) {
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
