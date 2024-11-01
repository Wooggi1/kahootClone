package service

import (
	"fmt"


	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/collection"
	"quiz.com/quiz/internal/entity"


)

type UserService struct {
	userCollection *collection.UserCollection
}

type LoginRequest struct {
	Email 		string
	Password 	string
}



func User(userCollection *collection.UserCollection) *UserService {
	return &UserService{
		userCollection: userCollection,
	}
}

func (s UserService) GetUsers() ([]entity.User, error) {
	return s.userCollection.GetAllUsers()
}

func (s UserService) GetUsersById(id primitive.ObjectID) (*entity.User, error) {
	return s.userCollection.GetUserById(id)
}

func (s UserService) HandleRegister(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	var newUser entity.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	filter := bson.M{"email": newUser.Email}
	exists, err := s.userCollection.CheckUserAlreadyExist(filter)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	newUser.Password = string(hashedPassword)

	err = s.userCollection.InsertUser(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error registering user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully!"})
}



func (s UserService) HandleLogin(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	var loginRequest LoginRequest

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	email:= loginRequest.Email
	user, err := s.userCollection.GetUserByEmail(email)

	if err != nil {
		fmt.Println("Error retrieving user")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "wrong password"})
	}

	token, err := s.userCollection.CreateToken(user)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error generating jwt key"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}