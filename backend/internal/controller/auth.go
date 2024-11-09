package controller

import (
	"github.com/gofiber/fiber/v2"
	"quiz.com/quiz/internal/service"
)

type UserController struct {
	userService *service.UserService
}

func User(userService *service.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (u *UserController) Register(c *fiber.Ctx) error {
	return u.userService.HandleRegister(c)
}

func (u *UserController) Login(c *fiber.Ctx) error {
	return u.userService.HandleLogin(c)
}

func (u *UserController) DetailUser(c *fiber.Ctx) error {
	return u .userService.DetailUser(c)
}

