package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend interface {
	SetAlive(bool)
	IsAlive() bool
	GetUrl() *url.URL
	Serve(http.ResponseWriter, *http.Request)
}

type backend struct {
	Url          *url.URL
	Alive        bool
	Mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
	Connections  int
}

// sets backend alive
func (b *backend) SetAlive(alive bool) {
	b.Mux.Lock()
	b.Alive = alive
	b.Mux.Unlock()
}

// returns current status of backend
func (b *backend) IsAlive() (alive bool) {
	b.Mux.RLock()
	defer b.Mux.RUnlock()
	alive = b.Alive
	return
}

func (b *backend) GetUrl() *url.URL {
	return b.Url
}

func (b *backend) Serve(rw http.ResponseWriter, req *http.Request) {
	b.ReverseProxy.ServeHTTP(rw, req)
}

func NewBackend(u *url.URL,rp *httputil.ReverseProxy) Backend{
	return &backend{
		Url : u,
		ReverseProxy: rp,
		Alive: true,
	}
}