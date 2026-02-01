package handlers

import "net/http"

type Router struct {
	handler http.Handler
	mux     *http.ServeMux
}

func NewRouter() *Router {
	router := &Router{
		mux: http.NewServeMux(),
	}

	router.registerRoutes()

	router.handler = router.mux

	return router
}

func (r *Router) registerRoutes() {
	r.mux.HandleFunc("/", Home)
}

func (r *Router) Handler() http.Handler {
	return r.handler
}
