package service

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"quiz.com/quiz/internal/entity"
	"quiz.com/quiz/internal/storage"
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

type LeaderboardEntry struct {
	Name 						string `json:"name"`
	Points					int	   `json:"points"`
}

type Game struct {
	Id              uuid.UUID
	Quiz            storage.Quiz
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

func newGame(quiz storage.Quiz, hostPlayer Player, host *websocket.Conn, netService *NetService) Game {
	fmt.Printf("New game created, host as player %s\n",hostPlayer.Name)
	hostPlayer.Connection = host
	hostPlayer.Points = 0
	hostPlayer.LastAwardedPoints = 0
	hostPlayer.Answered = false
	return Game{
		Id: uuid.New(),
		Quiz: quiz,
		CurrentQuestion: -1,
		Code: generateCode(),
		State: LobbyState,
		Ended: false,
		Time: 120,
		Players: []*Player{&hostPlayer},
		Host: host,
		netService: netService,
	}
}

func (g *Game) StartOrSkip() {
	if g.State == LobbyState {
		g.Start()
	} else {
		g.NextQuestion()
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
	g.Time = 120

	g.BroadcastPacket(QuestionShowPacket {
		Question: g.getCurrentQuestion(),
	}, true)
}

func (g *Game) Reveal() {
	g.Time = 5

	for _, player := range g.Players {
		g.netService.SendPacket(player.Connection, PlayerRevealPacket{
			Points: player.LastAwardedPoints,
		})
	}

	g.ChangeState(RevealState)
}

func (g *Game) Tick() {
	g.Time--
	g.BroadcastPacket(TickPacket{
		Tick: g.Time,
	}, true)

	if g.Time == 0 {
		switch g.State {
		case PlayState: 
			{
				g.Reveal()
				break
			}
		case RevealState:
			{
				g.Intermission()
				break
			}
		case IntermissionState:
			{
				g.NextQuestion()
				break
			}
		}
	}
}

func (g *Game) Intermission() {
	g.Time = 10
	g.ChangeState(IntermissionState)
	
	// Log the leaderboard data before broadcasting
	leaderboard := g.getLeaderboard()
	
	g.BroadcastPacket(LeaderboardPacket{
		Points: leaderboard,
	}, true)
}

func (g *Game) getLeaderboard() []LeaderboardEntry {
	// Sort players by points in descending order
	sort.Slice(g.Players, func(i, j int) bool {
		return g.Players[i].Points > g.Players[j].Points
	})

	leaderboard := []LeaderboardEntry{}
	for i := 0; i < int(math.Min(3, float64(len(g.Players)))); i++ {
		player := g.Players[i]
		leaderboard = append(leaderboard, LeaderboardEntry{
			Name: player.Name,
			Points: player.Points,
		})
	}
	
	return leaderboard
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