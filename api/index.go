package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/trevor-atlas/weekend/api/router"
	"github.com/trevor-atlas/weekend/api/users"
	"github.com/trevor-atlas/weekend/api/utils"
	"log"
	"net/http"
)

// Handler is the function that Now calls for every request
// Please note that only one function should be exported in this file!
func Handler(w http.ResponseWriter, r *http.Request) {
	defer utils.Track(utils.Runtime("handler"))
	log.Println("Request url: ", r.URL.Path)
	log.Println("Request method: ", r.Method)

	router.NewStupidRouter("/api").
		GET("/greet", greet).
		GET("/users", getUser).
		POST("/users", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(w, "not implemented yet!")
		}).
		Group("/v2", func(instance *router.StupidRouter) *router.StupidRouter {
			return instance.
				GET("/greet", greetV2).
				GET("/users", getUserV2)
		}).
		Group("/v3", func(instance *router.StupidRouter) *router.StupidRouter {
			return instance.POST("/greet", func(writer http.ResponseWriter, request *http.Request) {
				fmt.Fprintf(w, "this is a beta API! Be careful.")
			})
		}).
		DEFAULT(func (w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "I can't do that.")
		}).
		Start(w, r)

	log.Println("finished")
}

func greet(w http.ResponseWriter, r *http.Request) {
	name := utils.GetParamWithDefault("name", "Guest", r)
	fmt.Fprintf(w, `<center><h1>Hello, %s from Go on Now!</h1></center>`, name)
}

func greetV2(w http.ResponseWriter, r *http.Request) {
	name := utils.GetParamWithDefault("name", "Guest", r)
	fmt.Fprintf(w, `<center><h1>Hello, %s from APIV2 Go on Now!</h1></center>`, name)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	stub := &users.User{
		ID: uuid.New().String(),
		Name: "Vin Diesel",
	}
	utils.Write(stub, w)
}

func getUserV2(w http.ResponseWriter, r *http.Request) {
	stub := &users.User{
		ID: uuid.New().String(),
		Name: "John Hamm (v2)",
	}
	utils.Write(stub, w)
}
