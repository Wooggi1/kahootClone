package controller

import (
	"github.com/gofiber/fiber/v2"
	"quiz.com/quiz/internal/service"
)

type QuizController struct {
	quizService *service.QuizService
}

func Quiz(quizService *service.QuizService) QuizController {
	return QuizController{
		quizService: quizService,
	}
}

func (c QuizController) GetQuizzes(ctx *fiber.Ctx) error {
	quizzes, err := c.quizService.GetQuizzes()

	if err != nil {
		return err
	}

	return ctx.JSON(quizzes)
}

func (c QuizController) CreateQuiz(ctx *fiber.Ctx) error {
	return c.quizService.CreateQuiz(ctx)
}