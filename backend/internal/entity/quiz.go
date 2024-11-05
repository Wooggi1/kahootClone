package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizQuestions struct {
	Id      primitive.ObjectID	`json:"_id"`
	Name    string        			`json:"name"`
	Subject string							`json:"subject"`
	Content Content							`json:"content"`
	Choices []QuizChoices 			`json:"choices"`
}

type Content struct {
	Type		string	`json:"type"`
	Data    string	`json:"data"`
}

type QuizChoices struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}
