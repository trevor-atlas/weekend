package router

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type StupidRouter struct {
	base string
	routes map[string]http.HandlerFunc
	errors []error
	debug bool
}

func NewStupidRouter(basePath string) *StupidRouter {
	return &StupidRouter {
		base: trimLeadingSlash(basePath),
		routes: make(map[string]http.HandlerFunc),
		errors: make([]error, 0),
		// set this to false
		debug: getEnvAsBool("DEBUG", true),
	}
}

func (sr *StupidRouter) register(method string, route string, handler http.HandlerFunc) *StupidRouter {
	key := hashPath(method, sr.base, route)
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
		logMsg("Error building router:\n%s", message)
		http.Error(w, message, 500)
	}
	key := hashPath(r.Method, "", r.URL.Path)
	if cb, ok := sr.routes[key]; ok {
		logMsg("matched key: ", key)
		cb(w, r)
		return
	}
	if fallback, ok := sr.routes["DEFAULT"]; ok {
		logMsg("Use DEFAULT handler")
		fallback(w, r)
		return
	} else {
		logMsg("Use generic 404 handler")
		http.Error(w, "Not found", 404)
		return
	}
}

func hashPath(method, base, route string) string {
	return path.Join(method + ": ", base, route)
}

func trimLeadingSlash(str string) string {
	if str[len(str)-1:] == "/" {
		str = str[:len(str)-1]
	}
	return str
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

func logMsg(m ...interface{}) {
	if ok := getEnvAsBool("DEBUG", true); ok {
		log.Println(m)
	}
}
