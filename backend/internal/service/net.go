package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/entity"
	"quiz.com/quiz/internal/storage"
)

type NetService struct {
	quizService *QuizService
	games       []*Game
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		games:       []*Game{},
	}
}

type ConnectPacket struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type HostGamePacket struct {
	QuizId string `json:"quizId"`
  HostPlayer Player `json:"hostPlayer"`
}

type QuestionShowPacket struct {
	Question entity.QuizQuestions `json:"question"`
}

type ChangeGameStatePacket struct {
	State GameState `json:"state"`
}

type PlayerJoinPacket struct {
	Player Player `json:"player"`
}

type StartGamePacket struct{}

type TickPacket struct {
	Tick int `json:"tick"`
}

type QuestionAnswerPacket struct {
	Question int `json:"question"`
}

type PlayerRevealPacket struct {
	Points int `json:"points"`
}

type LeaderboardPacket struct {
	Points []LeaderboardEntry `json:"points"`
}

type GameCreatedPacket struct {
	Code string `json:"code"`
}

func (c *NetService) packetIdToPacket(packetId uint8) any {
	switch packetId {
	case 0:
		return &ConnectPacket{}
	case 1:
		return &HostGamePacket{}
	case 5:
		return &StartGamePacket{}
	case 7:
		return &QuestionAnswerPacket{}
	}
	return nil
}

func (c *NetService) packetToPacketId(packet any) (uint8, error) {
	switch packet.(type) {
	case QuestionShowPacket:
		return 2, nil
	case ChangeGameStatePacket:
		return 3, nil
	case PlayerJoinPacket:
		return 4, nil
	case TickPacket:
		return 6, nil
	case PlayerRevealPacket:
		return 8, nil
	case LeaderboardPacket:
		return 9, nil
	case GameCreatedPacket:
		return 10, nil
	}
	return 0, errors.New("invalid packet type")
}

func (c *NetService) getGameByCode(code string) *Game {
	for _, game := range c.games {
		if game.Code == code {
			return game
		}
	}
	return nil
}

func (c *NetService) getGameByHost(host *websocket.Conn) *Game {
	for _, game := range c.games {
		if game.Host == host {
			return game
		}
	}
	return nil
}

func (c *NetService) getGameByPlayer(con *websocket.Conn) (*Game, *Player) {
	for _, game := range c.games {
		for _, player := range game.Players {
			if player.Connection == con {
				return game, player
			}
		}
	}
	return nil, nil
}

func (c *NetService) OnIncomingMessage(con *websocket.Conn, mt int, msg []byte) {
	if len(msg) < 2 {
		return
	}

	packetId := msg[0]
	data := msg[1:]

	packet := c.packetIdToPacket(packetId)
	if packet == nil {
		return
	}

	err := json.Unmarshal(data, &packet)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch data := packet.(type) {
	case *ConnectPacket:
		game := c.getGameByCode(data.Code)
		if game == nil {
			return
		}
		game.OnPlayerJoin(data.Name, con)

	case *HostGamePacket:
		quizId, err := primitive.ObjectIDFromHex(data.QuizId)
		if err != nil {
			fmt.Println("Error parsing QuizId:", err)
			return
		}

		quiz, err := storage.GetQuizById(quizId)
		if err != nil {
			fmt.Println("Error fetching quiz:", err)
			return
		}

		if quiz == nil {
			fmt.Println("Quiz not found")
			return
		}

		newGame := newGame(*quiz, data.HostPlayer, con, c)
		fmt.Printf("New game created with code %s", newGame.Code)
		c.games = append(c.games, &newGame)

		c.SendPacket(con, GameCreatedPacket{
			Code: newGame.Code,
		})

		c.SendPacket(con, ChangeGameStatePacket{
			State: LobbyState,
		})

	case *StartGamePacket:
		game := c.getGameByHost(con)
		if game == nil {
			fmt.Println("game nulo - No game found for the provided host connection")
			return
		}

		fmt.Printf("Starting game with code: %s\n", game.Code)
		game.StartOrSkip()

	case *QuestionAnswerPacket:
		game, player := c.getGameByPlayer(con)
		if game == nil {
			return
		}
		game.OnPlayerAnswer(data.Question, player)
	}
}

func (c *NetService) SendPacket(connection *websocket.Conn, packet any) error {
	bytes, err := c.PacketToBytes(packet)
	if err != nil {
		return err
	}
	return connection.WriteMessage(websocket.BinaryMessage, bytes)
}

func (c *NetService) PacketToBytes(packet any) ([]byte, error) {
	packetId, err := c.packetToPacketId(packet)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	final := append([]byte{packetId}, bytes...)
	return final, nil
}
