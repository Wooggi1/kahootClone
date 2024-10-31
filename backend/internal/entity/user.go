package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        	primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name      	string							`json:"name"`
	Email     	string							`json:"email" bson:"email"`
	Password  	string							`json:"password"`
	TotalPoints int									`json:"points"`
}