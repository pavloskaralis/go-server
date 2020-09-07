package main 

import (
	"go-server/controller"
	"go-server/config/db"
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/kabukky/httpscerts"
	"github.com/joho/godotenv"
	"go-server/config/auth"
	"encoding/json"
)


 func Middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check expiration
		var resErr controller.ResponseError
		err := auth.ValidateToken(r)
		if err != nil {
			resErr.Error = "Token is invalid."
			json.NewEncoder(w).Encode(resErr)
			return
		}
        h.ServeHTTP(w, r)
    })
}



func main() {
	//load env variables (SIGNATURE)
	godotenv.Load()
	//note: self signed certificate; will cause Postman warning
    err := httpscerts.Check("cert.pem", "key.pem")
    if err != nil {
        err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8080")
        if err != nil {
            log.Fatal("Error: Could not create https certs.")
        }
	}

	//init mongo connection
	err = db.InitDB()
	if err != nil {
        log.Fatal("Error: Could not connect to mongoDB.")
	}

	//init Redis connection
	err = auth.InitRedis()
	if err != nil {
        log.Fatal("Error: Could not connect to redis.")
	}

	//router
	r := mux.NewRouter()
	r.HandleFunc("/signup", controller.SignupHandler).
		Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).
		Methods("POST")
	r.HandleFunc("/refresh", controller.RefreshHandler).
		Methods("POST")
	r.Handle("/profile", Middleware(http.HandlerFunc(controller.ProfileHandler))).
		Methods("GET")

	//connect
	fmt.Printf("listening on port 8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
}

