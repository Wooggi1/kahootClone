package collection

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"quiz.com/quiz/internal/entity"
)

type QuizCollection struct {
	collection *mongo.Collection
}

func Quiz(collection *mongo.Collection) *QuizCollection {
	return &QuizCollection{
		collection: collection,
	}
}

func (c QuizCollection) GetQuestionByAnySubject(filter bson.M) ([]entity.QuizQuestions ,error) {
	cursor, err := c.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var questions []entity.QuizQuestions
	err = cursor.All(context.Background(), &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}