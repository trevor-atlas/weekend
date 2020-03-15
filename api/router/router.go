package router

import (
	"errors"
	"log"
	"net/http"
)

type StupidRouter struct {
	base string
	routes map[string]http.HandlerFunc
	errors []error
}

func NewStupidRouter(basePath string) *StupidRouter {
	return &StupidRouter{
		base: basePath,
		routes: make(map[string]http.HandlerFunc),
		errors: make([]error, 0),
	}
}

func (sr *StupidRouter) register(method string, route string, handler http.HandlerFunc) *StupidRouter {
	key := trimLeadingSlash(method + "-" + sr.base + route)
	if _, ok := sr.routes[key]; ok {
		err := errors.New("More than one route matched <" + route + "> for method <" + method + ">")
		sr.errors = append(sr.errors, err)
	}
	sr.routes[key] = handler
	return sr
}

func (sr *StupidRouter) Group(basePath string, fn func(instance *StupidRouter) *StupidRouter) *StupidRouter {
	temp := fn(NewStupidRouter(sr.base + basePath))
	for k, v := range temp.routes {
		sr.routes[k] = v
	}
	sr.errors = append(sr.errors, temp.errors...)
	return sr
}

func (sr *StupidRouter) GET(path string, handler http.HandlerFunc) *StupidRouter {
	return sr.register(http.MethodGet, path, handler)
}

func (sr *StupidRouter) POST(path string, handler http.HandlerFunc) *StupidRouter {
	return sr.register(http.MethodPost, path, handler)
}

func (sr *StupidRouter) PUT(path string, handler http.HandlerFunc) *StupidRouter {
	return sr.register(http.MethodPut, path, handler)
}

func (sr *StupidRouter) DEFAULT(errorHandler http.HandlerFunc) *StupidRouter {
	sr.routes["DEFAULT"] = errorHandler
	return sr
}

func (sr *StupidRouter) Start(w http.ResponseWriter, r *http.Request) {
	if len(sr.errors) > 0 {
		var message string
		for _, err := range sr.errors {
			message += err.Error() + ",\n"
		}
		http.Error(w, message, 500)
	}
	key := trimLeadingSlash(r.Method + "-" + r.URL.Path)
	log.Println("matched key: ", key)
	if cb, ok := sr.routes[key]; ok {
		cb(w, r)
		return
	}
	if fallback, ok := sr.routes["DEFAULT"]; ok {
		fallback(w, r)
		return
	} else {
		http.Error(w, "Not found", 404)
		return
	}
}

func trimLeadingSlash(str string) string {
	if str[len(str)-1:] == "/" {
		str = str[:len(str)-1]
	}
	return str
}
