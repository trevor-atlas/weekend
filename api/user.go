package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/trevor-atlas/weekend/api/users"
	"log"
	"net/http"
	"github.com/trevor-atlas/weekend/api/utils"
)

// Handler is the function that Now calls for every request
func User(w http.ResponseWriter, r *http.Request) {
	defer utils.Track(utils.Runtime("handler"))
	log.Println("Request url: ", r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		getUser(w)
	case http.MethodPost:
		users.CreateUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func getUser(w http.ResponseWriter) {
	stub := &users.User {
		ID: uuid.New().String(),
		Name: "Alex Jones",
		Token: "lol yeah right",
	}
	utils.Write(stub, w)
}
