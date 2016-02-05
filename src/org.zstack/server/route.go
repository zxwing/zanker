package server

import (
	"fmt"
	LOG "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	ApiRoute interface {
		Methods() []string
		Path() []string
		Handler(http.ResponseWriter, *http.Request)
	}

	RouteManager struct{}
)

var (
	routes []ApiRoute = make([]ApiRoute, 0)

	routeManager = &RouteManager{}
)

func RegisterApiRoute(r ApiRoute) {
	paths := r.Path()

	if len(paths) == 0 {
		panic("Path cannot be empty")
	}

	for _, p := range paths {
		if p == "" {
			panic("empty path is not allowed")
		}
	}

	if len(r.Methods()) == 0 {
		panic("Methods cannot be empty")
	}

	for _, p := range paths {
		LOG.WithFields(LOG.Fields{
			"Path":    p,
			"Methods": r.Methods(),
		}).Debug("Registered a new Handler")
	}

	routes = append(routes, r)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		for _, path := range r.Path() {
			re := router.HandleFunc(Url(path).Path(), routeManager.WrapHandler(r))
			re.Methods(r.Methods()...)
		}
	}

	return router
}

func (mgr *RouteManager) WrapHandler(api ApiRoute) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
			}
		}()

		LOG.Debugf("%s %v", req.Method, req.URL)
		api.Handler(w, req)
	}
}
