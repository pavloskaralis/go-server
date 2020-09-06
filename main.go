package main 

import (
	"go-server/controller"
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/kabukky/httpscerts"
	"github.com/joho/godotenv"
)


func main() {
	godotenv.Load()
	//note: self signed certificate
    err := httpscerts.Check("cert.pem", "key.pem")
    if err != nil {
        err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8080")
        if err != nil {
            log.Fatal("Error: Couldn't create https certs.")
        }
	}

	r := mux.NewRouter()
	r.HandleFunc("/signup", controller.SignupHandler).
		Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).
		Methods("POST")
	r.HandleFunc("/profile", controller.ProfileHandler).
		Methods("GET")

	fmt.Printf("listening on port 8080")

	http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r)
	
}

