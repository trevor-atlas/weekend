package router

import (
	"errors"
	"log"
	"net/http"
)

type stupidRouter struct {
	routes map[string]http.HandlerFunc
	errors []error
}

func NewStupidRouter() *stupidRouter {
	return &stupidRouter{
		routes: make(map[string]http.HandlerFunc),
		errors: make([]error, 0),
	}
}

func (sr *stupidRouter) register(method string, route string, handler http.HandlerFunc) *stupidRouter {
	key := method + "-" + route
	if _, ok := sr.routes[key]; ok {
		err := errors.New("More than one route matched <" + route + "> for method <" + method + ">")
		sr.errors = append(sr.errors, err)
	}
	sr.routes[key] = handler
	return sr
}

func (sr *stupidRouter) GET(path string, handler http.HandlerFunc) *stupidRouter {
	return sr.register(http.MethodGet, path, handler)
}

func (sr *stupidRouter) POST(path string, handler http.HandlerFunc) *stupidRouter {
	return sr.register(http.MethodPost, path, handler)
}

func (sr *stupidRouter) PUT(path string, handler http.HandlerFunc) *stupidRouter {
	return sr.register(http.MethodPut, path, handler)
}

func (sr *stupidRouter) DEFAULT(errorHandler http.HandlerFunc) *stupidRouter {
	sr.routes["DEFAULT"] = errorHandler
	return sr
}

func (sr *stupidRouter) Start(w http.ResponseWriter, r *http.Request) {
	if len(sr.errors) > 0 {
		var message string
		for _, err := range sr.errors {
			message += err.Error() + ",\n"
		}
		http.Error(w, message, 500)
	}
	key := r.Method + "-" + r.URL.Path
	log.Println("matched key: ", key)
	// strip leading slash /
	if key[len(key)-1:] == "/" {
		key = key[:len(key)-1]
	}
	log.Println("check map key: ", key)
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

