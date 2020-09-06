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

