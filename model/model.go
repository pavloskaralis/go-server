package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	UID primitive.ObjectID `bson:"_id,omitempty" json:"uid,omitempty"`
}

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