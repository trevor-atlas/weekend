package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/trevor-atlas/weekend/api/users"
	"github.com/trevor-atlas/weekend/api/utils"
	"log"
	"net/http"
)

// Handler is the function that Now calls for every request
func Handler(w http.ResponseWriter, r *http.Request) {
	defer utils.Track(utils.Runtime("handler"))
	log.Println("Request url: ", r.URL.Path)

	switch r.URL.Path {
	case "/api/greet":
	case "/api/greet/":
		greet(w, r)
	case "/api/users":
	case "/api/users/":
		switch r.Method {
		case http.MethodGet:
			getUser(w)
		case http.MethodPost:
			users.CreateUser(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "I can't do that.")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	name := utils.GetParamWithDefault("name", "Guest", r)
	//w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<center><h1>Hello, %s from Go on Now!</h1></center>`, name)
}

func getUser(w http.ResponseWriter) {
	stub := &users.User{
		ID: uuid.New().String(),
		Name: "Alex Jones",
		Token: "lol yeah right",
	}
	utils.Write(stub, w)
}
