package router

import (
	"github.com/gorilla/mux"
	logger "image_gallery/logger"
	"net/http"
)

// DefaultRouteScheme is the default route scheme
const DefaultRouteScheme string = "http"

// Router is an http router
type Router struct {
	Handlers    []Handler
	Logger      *logger.Logger
}

// Route struct defining all routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Scheme      string
	HandlerFunc http.HandlerFunc
}

// Routes slice of Route
type Routes []Route

// Handler is an handler that exposes routes
type Handler interface {
	Routes() Routes
}

// AddHandler adds a new handler
func (r *Router) AddHandler(h Handler) {
	if r.Handlers == nil {
		r.Handlers = make([]Handler, 0)
	}

	r.Handlers = append(r.Handlers, h)
}

// Configure registers all handlers routes and return mux router
func (r *Router) Configure() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, handler := range r.Handlers {
		for _, route := range handler.Routes() {
			handlerFunc := route.HandlerFunc
			if route.Scheme == "" {
				route.Scheme = DefaultRouteScheme
			}
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Schemes(route.Scheme).
				Handler(handlerFunc)
		}
	}

	return router
}
