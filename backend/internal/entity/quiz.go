package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Questions []QuizQuestions    `json:"questions"`
}

type QuizQuestions struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Subject string				`json:"subject"`
	Choices []QuizChoices `json:"choices"`
}

type QuizChoices struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}
