package service

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"quiz.com/quiz/internal/entity"
)

type Player struct {
	Id 								uuid.UUID					`json:"id"`
	Name							string						`json:"name"`
	Connection				*websocket.Conn		`json:"-"`
	Points						int								`json:"-"`
	LastAwardedPoints	int								`json:"-"`
	Answered					bool							`json:"-"`
}

type GameState int 

const (
	LobbyState GameState = iota
	PlayState
	IntermissionState
	RevealState
	EndState
)

type Game struct {
	Id              uuid.UUID
	Quiz            entity.Quiz
	CurrentQuestion int
	Code            string
	State 					GameState
	Ended 					bool
	Time						int
	Players					[]*Player
	

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
		CurrentQuestion: -1,
		Code: generateCode(),
		State: LobbyState,
		Ended: false,
		Time: 60,
		Players: []*Player{},
		Host: host,
		netService: netService,
	}
}

func (g *Game) Start() {
	g.ChangeState(PlayState)
	g.NextQuestion()
	
	go func() {
		for {
			if(g.Ended) {
				return
			}

			g.Tick()
			time.Sleep(time.Second)
		}
	} ()
}

func (g *Game) ResetPlayerAnswerStates() {
	for _, player := range g.Players {
		player.Answered = false
	}
}

func (g *Game) End() {
	g.Ended = true
	g.ChangeState(EndState)
}

func (g *Game) NextQuestion() {
	g.CurrentQuestion++

	if g.CurrentQuestion >= len(g.Quiz.Questions) {
		g.End()
		return
	}

	g.ChangeState(PlayState)
	g.Time = 60

	g.netService.SendPacket(g.Host, QuestionShowPacket {
		Question: g.getCurrentQuestion(),
	})
}

func (g *Game) Reveal() {
	g.Time = 10

	for _, player := range g.Players {
		g.netService.SendPacket(player.Connection, PlayerRevealPacket{
			Points: player.LastAwardedPoints,
		})
	}

	g.ChangeState(RevealState)
}

func (g *Game) Tick() {
	g.Time--
	g.netService.SendPacket(g.Host, TickPacket{
		Tick: g.Time,
	})

	if g.Time == 0 {
		switch g.State {
		case PlayState: {
			g.Reveal()
			break
		}
		case RevealState:{
			g.NextQuestion()
			break
		}
		}
	}
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

	g.Players = append(g.Players, &player)

	g.netService.SendPacket(connection, ChangeGameStatePacket{
		State: g.State,
	})

	g.netService.SendPacket(g.Host, PlayerJoinPacket{
		Player: player,
	})
}

func (g *Game) getAnsweredPlayers() []*Player {
	players := []*Player{}

	for _, player := range g.Players {
		if player.Answered {
			players = append(players, player)
		}
	}

	return players
}

func (g *Game) getCurrentQuestion() entity.QuizQuestions {
	return g.Quiz.Questions[g.CurrentQuestion]
}

func (g *Game) isCorrectChoice(choiceIndex int) bool {
	choices := g.getCurrentQuestion().Choices
	if choiceIndex < 0 || choiceIndex >= len(choices) {
		return false
	}

	return choices[choiceIndex].Correct
}

func (g *Game) getPointsReward() int {
	answered := len(g.getAnsweredPlayers())
	orderReward := 5000 - (1000 * math.Min(4, float64(answered)))
	timeReward := g.Time * (1000 / 60)

	return int(orderReward) + timeReward
}

func (g *Game) OnPlayerAnswer(choice int, player *Player) {
	if g.isCorrectChoice(choice) {
		player.LastAwardedPoints = g.getPointsReward()
		player.Points +=player.LastAwardedPoints
	} else {
		player.LastAwardedPoints = 0
	}

	player.Answered = true

	if len(g.getAnsweredPlayers()) == len(g.Players) {
		g.Reveal()
	}
}