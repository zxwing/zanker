package server

import (
	"fmt"
	LOG "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	Route struct {
		Methods []string
		Path    string
		Handler http.HandlerFunc
	}
)

var (
	routes []*Route = make([]*Route, 0)
)

func (r *Route) Register() {
	if r.Path == "" {
		panic("Path cannot be empty")
	}

	if r.Handler == nil {
		panic("Handler cannot be nil")
	}

	if len(r.Methods) == 0 {
		r.Methods = []string{"POST", "GET"}
	}

	routes = append(routes, r)
	LOG.WithFields(LOG.Fields{
		"Path":    r.Path,
		"Methods": r.Methods,
	}).Debug("Registered a new Handler")
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		re := router.HandleFunc(r.Path, r.GetWrappedHandler())
		re.Methods(r.Methods...)
	}

	return router
}

func (r *Route) GetWrappedHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
			}
		}()

		LOG.Debugf("%s %v", req.Method, req.URL)
		r.Handler(w, req)
	}
}
