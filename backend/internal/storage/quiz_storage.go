package storage

import (
  "quiz.com/quiz/internal/entity"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "errors"
)

type Quiz struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Questions []entity.QuizQuestions    `json:"questions"`
}

var Quizzes = make([]Quiz, 0)

func AddQuiz(quiz Quiz) {
  Quizzes = append(Quizzes, quiz)
}

func GetQuizzes() []Quiz {
  return Quizzes
}

func GetQuizById(id primitive.ObjectID) (*Quiz, error) {
  for _, quiz := range Quizzes {
    if quiz.Id == id {
      return &quiz, nil
    }
  }
  
	return nil, errors.New("quiz not found")
}
