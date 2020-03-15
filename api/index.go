package handler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/trevor-atlas/weekend/api/users"
	"github.com/trevor-atlas/weekend/api/utils"
	"log"
	"net/http"
)

type StupidRouter struct {
	routes map[string]http.HandlerFunc
	errors []error
}

func NewStupidRouter() *StupidRouter {
	sr := &StupidRouter{
		routes: make(map[string]http.HandlerFunc),
		errors: make([]error, 0),
	}
	return sr
}

func (sr *StupidRouter) register(method string, route string, handler http.HandlerFunc) {
	key := method + "-" + route
	if _, ok := sr.routes[key]; ok {
		err := errors.New("More than one route matched <" + route + "> for method <" + method + ">")
		sr.errors = append(sr.errors, err)
	}
	sr.routes[key] = handler
}

func (sr *StupidRouter) GET(path string, handler http.HandlerFunc) *StupidRouter {
	sr.register(http.MethodGet, path, handler)
	return sr
}

func (sr *StupidRouter) POST(path string, handler http.HandlerFunc) *StupidRouter {
	sr.register(http.MethodPost, path, handler)
	return sr
}

func (sr *StupidRouter) PUT(path string, handler http.HandlerFunc) *StupidRouter {
	sr.register(http.MethodPut, path, handler)
	return sr
}

func (sr *StupidRouter) DEFAULT(errorHandler http.HandlerFunc) *StupidRouter {
	sr.routes["DEFAULT"] = errorHandler
	return sr
}

func (sr *StupidRouter) Run(w http.ResponseWriter, r *http.Request) {
	if len(sr.errors) > 0 {
		var message string
		for _, err := range  sr.errors {
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

// Handler is the function that Now calls for every request
func Handler(w http.ResponseWriter, r *http.Request) {
	defer utils.Track(utils.Runtime("handler"))
	log.Println("Request url: ", r.URL.Path)
	log.Println("Request method: ", r.Method)

	NewStupidRouter().
		GET("/api/greet", greet).
		GET("/api/users", getUser).
		POST("/api/users", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(w, "not implemented yet!")
		}).
		DEFAULT(func (w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "I can't do that.")
		}).
		Run(w, r)
}

func greet(w http.ResponseWriter, r *http.Request) {
	name := utils.GetParamWithDefault("name", "Guest", r)
	fmt.Fprintf(w, `<center><h1>Hello, %s from Go on Now!</h1></center>`, name)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	stub := &users.User{
		ID: uuid.New().String(),
		Name: "Alex Jones",
		Token: "lol yeah right",
	}
	utils.Write(stub, w)
}
