package service

import (
	"github.com/gofiber/fiber/v2"
	"quiz.com/quiz/internal/collection"
	"quiz.com/quiz/internal/entity"
)

type QuizService struct {
	quizCollection *collection.QuizCollection
}

func Quiz(quizCollection *collection.QuizCollection) *QuizService {
	return &QuizService{
		quizCollection: quizCollection,
	}
}

func (s QuizService) GetQuizzes() ([]entity.Quiz, error) {
	return s.quizCollection.GetQuizzes()
}

func (s QuizService) CreateQuiz(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	var newQuiz entity.Quiz
	if err := c.BodyParser(&newQuiz); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	err := s.quizCollection.InsertQuiz(newQuiz)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error registering quiz"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Quiz created successfully!"})
}