package models

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
		Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		FirstName    string 					`json:"firstname"`
		LastName     string        		`json:"lastname"`
		Email        string        		`json:"email"`
		Password     string        		`json:"password,omitempty"`
		HashPassword []byte        		`json:"hashpassword,omitempty"`
	}
