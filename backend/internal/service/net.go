package service

import (
	"fmt"
	"strings"

	"github.com/gofiber/contrib/websocket"
)

type NetService struct {
	quizService *QuizService
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
	}
}

func (c *NetService) OnIncomingMessage(con *websocket.Conn, mt int, msg []byte) {
	str := string(msg)
	parts := strings.Split(str, ":")
	cmd := parts[0]
	argument := parts[1]

	switch cmd {
	case "host":
		{
			fmt.Println("host quiz", argument)
		}

	case "join":
		{
			fmt.Println("join code", argument)
		}
	}
}
