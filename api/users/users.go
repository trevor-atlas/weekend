package users

import (
	"fmt"
	"net/http"
)

type User struct {
	ID    string `fauna:"id"`
	Token string `fauna:"token"`
	Name  string `fauna:"name"`
}

type LoginRequest struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type UserCreateRequest struct {
	Name  string `form:"name" json:"name" xml:"name"  binding:"required"`
	Token string `form:"token" json:"name" xml:"name"  binding:"required"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login is unimplemented")
}
