package collection

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"quiz.com/quiz/internal/entity"
)

type UserCollection struct {
	collection *mongo.Collection
}

func User(collection *mongo.Collection) *UserCollection {
	return &UserCollection{
		collection: collection,
	}
}

func (c UserCollection) InsertUser(user entity.User) error {
	_, err := c.collection.InsertOne(context.Background(), user)
	return err
}

func (c UserCollection) GetAllUsers() ([]entity.User, error) {
	cursor, err := c.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var users []entity.User
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c UserCollection) GetUserById(id primitive.ObjectID) (*entity.User, error) {
	result := c.collection.FindOne(context.Background(), bson.M{"_id": id})

	var user entity.User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c UserCollection) CheckUserAlreadyExist(filter bson.M) (bool, error) {
	result := c.collection.FindOne(context.Background(), filter)

	var user entity.User
	err := result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}