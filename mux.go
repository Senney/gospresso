package gospresso

import (
	"net/http"
	"sync"
)

type Mux struct {
	handler http.Handler
	pool    *sync.Pool
	routes  *RouteTree
}

func NewMux() *Mux {
	mux := &Mux{pool: &sync.Pool{}, routes: NewRouteTree()}

	// todo: initialize mux.pool.New

	return mux
}

func (mx *Mux) Get(pattern string, handlerFn http.HandlerFunc) {
	mx.handle(mGET, pattern, handlerFn)
}

func (mx *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mx.handler == nil {
		panic("Handle no routes")
	}

	mx.handler.ServeHTTP(w, r)
}

func (mx *Mux) handle(method uint, pattern string, handlerFn http.HandlerFunc) {
	if mx.handler == nil {
		mx.handler = http.HandlerFunc(mx.routeHTTP)
	}

	mx.routes.root.Insert(method, pattern, handlerFn)
}

func (mx *Mux) routeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "" {
		path = "/"
	}

	route := mx.routes.root.Search(mGET, path)

	if route == nil {
		http.NotFound(res, req)
		return
	}

	route.Handler.ServeHTTP(res, req)
}
