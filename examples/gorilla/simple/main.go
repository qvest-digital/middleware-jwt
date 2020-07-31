package main

import (
	"net/http"

	"github.com/gorilla/mux"
	jwt "github.com/tarent/middleware-jwt/internal"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", YourHandler)
	r.Use(jwt.JwtAuthAnyGroup("mysecret", "groupA", "groupB"))

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Success!\n"))
}
