package server

import (
	"net/http"
)

type (
	_ping struct{}

	Ping _ping
)

var (
	PluginPing = &Ping{}
)

func (p *Ping) Methods() []string {
	return []string{"GET", "POST"}
}

func (p *Ping) Path() string {
	return "/ping"
}

func (p *Ping) Handler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func init() {
	RegisterApiRoute(PluginPing)
}
