package serverpool

import (
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/rimo10/load_balancer/backend"
)

type Serverpool interface {
	AddBackend(backend.Backend)
	GetBackends() []backend.Backend
	NextIndex() int
	GetNextPeerRoundRobin() backend.Backend
	MarkBackendStatus(*url.URL, bool)
}

type serverpool struct {
	backends []backend.Backend
	mux      sync.RWMutex
	current  uint64
}

func (s *serverpool) AddBackend(b backend.Backend) {
	s.backends = append(s.backends, b)
}

func (s *serverpool) GetBackends() []backend.Backend {
	return s.backends
}

func (s *serverpool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

func (s *serverpool) GetNextPeerRoundRobin() backend.Backend {
	next := s.NextIndex()
	l := next + len(s.backends)
	for idx := next; idx < l; idx++ {
		idx := idx % len(s.backends)
		if s.backends[idx].IsAlive() {
			if idx != next {
				atomic.AddUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

func (s *serverpool) MarkBackendStatus(u *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.GetUrl().String() == u.String() {
			b.SetAlive(alive)
			break
		}
	}
}

func NewServerPool() (Serverpool, error) {
	return &serverpool{
		backends: make([]backend.Backend, 0),
		current:  0,
	}, nil
}
