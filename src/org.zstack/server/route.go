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
		Path() string
		Handler(http.ResponseWriter, *http.Request)
	}

	RouteManager struct{}
)

var (
	routes []ApiRoute = make([]ApiRoute, 0)

	routeManager = &RouteManager{}
)

func RegisterApiRoute(r ApiRoute) {
	if r.Path() == "" {
		panic("Path cannot be empty")
	}

	if len(r.Methods()) == 0 {
		panic("Methods cannot be empty")
	}

	LOG.WithFields(LOG.Fields{
		"Path":    r.Path(),
		"Methods": r.Methods(),
	}).Debug("Registered a new Handler")

	routes = append(routes, r)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		re := router.HandleFunc(Url(r.Path()).Path(), routeManager.WrapHandler(r))
		re.Methods(r.Methods()...)
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
