package collection

import (
	"context"
	"errors"
	
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"quiz.com/quiz/internal/entity"
	"github.com/golang-jwt/jwt/v5"

)


type UserCollection struct {
	collection *mongo.Collection
	secret_key []byte
}

func User(collection *mongo.Collection, secret_key []byte) *UserCollection {
	return &UserCollection{
		collection: collection,
		secret_key: secret_key,
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

func (c UserCollection) GetUserByEmail(email string) (*entity.User, error) {
	result := c.collection.FindOne(context.Background(), bson.M{"email": email})

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

func (c UserCollection) CreateToken(user *entity.User) (string, error) {
	// Check if the secret key is set
	if len(c.secret_key) == 0 {
		return "", errors.New("secret key not set")
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(), // Token expires in 30 days
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(c.secret_key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}