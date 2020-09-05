package main 

import (
	"go-server/controller"
	"go-server/model"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)



func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", controller.signupHandler).
		Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).
		Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}