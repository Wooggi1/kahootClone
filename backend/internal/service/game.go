package service

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"quiz.com/quiz/internal/entity"
)

type Player struct {
	Id 							uuid.UUID					`json:"id"`
	Name						string						`json:"name"`
	Connection			*websocket.Conn		`json:"-"`
}

type GameState int 

const (
	LobbyState GameState = iota
	PlayState
	RevealState
	EndState
)

type Game struct {
	Id              uuid.UUID
	Quiz            entity.Quiz
	Code            string
	State 					GameState
	Time						int
	Players					[]Player

	Host 						*websocket.Conn
	netService			*NetService
}

func generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

func newGame(quiz entity.Quiz, host *websocket.Conn, netService *NetService) Game {
	return Game{
		Id: uuid.New(),
		Quiz: quiz,
		Code: generateCode(),
		State: LobbyState,
		Time: 60,
		Players: []Player{},
		Host: host,
		netService: netService,
	}
}

func (g *Game) Start() {
	g.ChangeState(PlayState)
	g.netService.SendPacket(g.Host, QuestionShowPacket{
		Question: entity.QuizQuestions{
			Id: "",
			Name: "What is 2 + 2?",
			Choices: []entity.QuizChoices{
				{
					Id: 	"a",
					Name: "4",	
				},
				{
					Id: 	"b",
					Name: "5",	
				},
				{
					Id: 	"c",
					Name: "6",	
				},
				{
					Id: 	"d",
					Name: "7",	
				}},
		},
	})
	
	go func() {
		for {
			g.Tick()
			time.Sleep(time.Second)
		}
	} ()
}

func (g *Game) Tick() {
	g.Time--
	g.netService.SendPacket(g.Host, TickPacket{
		Tick: g.Time,
	})
}

func (g *Game) ChangeState(state GameState) {
	g.State = state 
	g.BroadcastPacket(ChangeGameStatePacket{
		State: state,
	}, true)
}

func (g *Game) BroadcastPacket(packet any, includeHost bool) error {
	for _, player := range g.Players {
		err := g.netService.SendPacket(player.Connection, packet)
		if err != nil {
			return err
		}
	}

	if includeHost{
		err := g.netService.SendPacket(g.Host, packet)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) OnPlayerJoin(name string, connection *websocket.Conn) {
	fmt.Println(name, "joined the game")

	player := Player{
		Id:					uuid.New(),	
		Name: 			name,
		Connection: connection,
	}

	g.Players = append(g.Players, player)

	g.netService.SendPacket(connection, ChangeGameStatePacket{
		State: g.State,
	})

	g.netService.SendPacket(g.Host, PlayerJoinPacket{
		Player: player,
	})
}