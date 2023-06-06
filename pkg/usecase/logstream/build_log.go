package logstream

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
