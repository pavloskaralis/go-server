package main 

import (
	"go-server/controller"
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Printf("listening on port 8080")

	r := mux.NewRouter()
	r.HandleFunc("/signup", controller.SignupHandler).
		Methods("POST")
	// r.HandleFunc("/login", controller.LoginHandler).
	// 	Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}