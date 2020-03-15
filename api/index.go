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
func Handler(w http.ResponseWriter, r *http.Request) {
	defer utils.Track(utils.Runtime("handler"))
	log.Println("Request url: ", r.URL.Path)
	log.Println("Request method: ", r.Method)

	yolo := router.NewStupidRouter().
		GET("/api/greet", greet).
		GET("/api/users", getUser).
		POST("/api/users", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(w, "not implemented yet!")
		}).
		DEFAULT(func (w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "I can't do that.")
		})

		yolo.Start(w, r)

	log.Println("finished")
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
