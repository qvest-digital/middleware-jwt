package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jwt "github.com/tarent/middleware-jwt/internal"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", yourHandler)
	r.Use(jwt.JwtAuthAllowAll("mysecret"))

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func yourHandler(w http.ResponseWriter, r *http.Request) {

	groups := jwt.GetGroupsFromAuthenticatedRequest(r)

	response := fmt.Sprintf("My groups are: %v\n", groups)

	w.Write([]byte(response))
}
