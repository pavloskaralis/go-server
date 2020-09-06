package controller

type ResponseError struct {
	Error string `json:"error"`
}

type Auth struct {
	Access string `json:"access"`
	Refresh string `json:"refresh"`
} 

type Profile struct {
	UID string `json:"uid"`
	Username string `json:"username"`
	Email string `json:"email"`
} 

type ResponseSuccess struct {
	Auth Auth
	Profile Profile 
}