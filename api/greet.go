package handler

import (
	"fmt"
	"log"
	"net/http"
	u "github.com/trevor-atlas/weekend/api/utils"
)

// Handler is the function that Now calls for every request
func Greet(w http.ResponseWriter, r *http.Request) {
	defer u.Track(u.Runtime("handler"))
	log.Println("Request url: ", r.URL.Path)

	name := u.GetParamWithDefault("name", "Guest", r)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<center><h1>Hello, %s from Go on Now!</h1></center>`, name)
}
