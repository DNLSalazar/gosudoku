package server

import (
	"errors"
	"fmt"

	"github.com/DNLSalazar/gosudoku/sudoku"
)

type MessageFromPlayer struct {
	Coors    sudoku.Coor `json:"coors"`
	Value    int         `json:"value"`
	PlayerId string      `json:"playerId"`
}

type PlayerSession struct {
	PlayerId   string
	PlayerName string
	Game       sudoku.Sudoku
}

func (p *PlayerSession) ToPlayerSessionResponse() PlayerGameStartResponse {
	return PlayerGameStartResponse{
		PlayerId:   p.PlayerId,
		Board:      p.Game.GetBoard(),
		PlayerName: p.PlayerName,
	}
}

func (p *PlayerSession) ToOpponenetInformation() OpponentPlayerInformation {
	return OpponentPlayerInformation{
		PlayerName: p.PlayerName,
		Board:      p.Game.GetDumbBoard(),
	}
}

type PlayerGameStartResponse struct {
	PlayerId   string          `json:"playerId"`
	PlayerName string          `json:"playerName"`
	Board      [][]sudoku.Cell `json:"board"`
}

type OpponentPlayerInformation struct {
	PlayerName string              `json:"playerName"`
	Board      [][]sudoku.DumbCell `json:"board"`
}

type GameSession struct {
	InitialBoard [][]sudoku.Cell
	Players      []*PlayerSession
	MaxPlayers   int
	Winner       string
	onGoing      bool
}

type GameEndedInformation struct {
	WinnerName  string                      `json:"winnerName"`
	PlayersData []OpponentPlayerInformation `json:"playerData"`
}

type GameErrorInformation struct {
	Message string `json:"message"`
}

type MessageResponse[T [][]sudoku.Cell | []OpponentPlayerInformation | PlayerGameStartResponse | GameEndedInformation | GameErrorInformation] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

func (gs *GameSession) SlotAvailable() bool {
	return len(gs.Players) < gs.MaxPlayers
}

func (gs *GameSession) GetPlayerSession(id string) *PlayerSession {
	for _, v := range gs.Players {
		if id == v.PlayerId {
			return v
		}
	}
	return nil
}

func (gs *GameSession) AddPlayerSession(ps *PlayerSession) error {
	if !gs.SlotAvailable() {
		return errors.New("No slots for player session")
	}

	for _, v := range gs.Players {
		if ps.PlayerId == v.PlayerId {
			return errors.New("Player already in game")
		}
	}

	ps.PlayerName = fmt.Sprintf("Player%d", len(gs.Players)+1)
	gs.Players = append(gs.Players, ps)

	return nil
}

func (gs *GameSession) GetPlayersProgressInformation() []OpponentPlayerInformation {
	psInfo := make([]OpponentPlayerInformation, len(gs.Players))
	for i, v := range gs.Players {
		psInfo[i] = v.ToOpponenetInformation()
	}

	return psInfo
}
