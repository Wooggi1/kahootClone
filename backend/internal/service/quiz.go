package service

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/collection"
	"quiz.com/quiz/internal/entity"
	"quiz.com/quiz/internal/storage"
)

type QuizService struct {
	quizCollection *collection.QuizCollection
}

type QuizRequest struct {
	Name			string		`json:"name"`
	Subjects	[]string 	`json:"subjects"`
}

type GetQuizRequest struct {
	Id        string `json:"id"`
}

func Quiz(quizCollection *collection.QuizCollection) *QuizService {
	return &QuizService{
		quizCollection: quizCollection,
	}
}

func (s QuizService) CreateQuiz(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	
	var quizRequest QuizRequest

	if err := c.BodyParser(&quizRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if len(quizRequest.Subjects) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "select subjects"})
	}
	filter := bson.M{"subject":bson.M{"$in": quizRequest.Subjects}}

	questions, err := s.quizCollection.GetQuestionByAnySubject(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "couldnt retrieve questions"})
	}

	var quizQuestions []entity.QuizQuestions
	for _, question := range questions {
		quizQuestions = append(quizQuestions, entity.QuizQuestions{
			Id: primitive.NewObjectID(),
			Name: question.Name,
			Subject: question.Subject,
			Content: question.Content,
			Choices: question.Choices,
		})
	}

	newQuiz := storage.Quiz{
		Id: primitive.NewObjectID(),
		Name: quizRequest.Name,
		Questions: quizQuestions,
	}

	storage.AddQuiz(newQuiz)

	return c.Status(fiber.StatusCreated).JSON(newQuiz)
}