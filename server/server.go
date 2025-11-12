package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sudoku/sudoku"
	"time"

	"github.com/olahol/melody"
)

const PORT = "5000"

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

func Server() {
	m := melody.New()
	// game := sudoku.CreateNewSudoku(17)
	game := sudoku.CreateNewSudoku(76)
	gs := GameSession{
		InitialBoard: game.GetBoard(),
		MaxPlayers:   4,
		Players:      []*PlayerSession{},
		onGoing:      true,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "server/index.html")
	})

	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		if !gs.SlotAvailable() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Game session full"))
			return
		}

		playerId := fmt.Sprintf("%d", time.Now().UnixMilli())
		pSudoku := sudoku.CreateSudokuFromCells(game.GetBoard())
		pSession := PlayerSession{
			PlayerId: playerId,
			Game:     pSudoku,
		}

		err := gs.AddPlayerSession(&pSession)
		if err != nil {
			fmt.Println("Error adding player session")
			w.Write([]byte("Cannot add player session to game"))
			return
		}

		resData := MessageResponse[PlayerGameStartResponse]{
			Type: "session",
			Data: pSession.ToPlayerSessionResponse(),
		}

		data, err := json.Marshal(resData)
		if err != nil {
			fmt.Println("Error getting sudoku")
			w.Write([]byte("Error getting sudoku"))
			return
		}

		w.Write(data)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if !gs.onGoing {
			errorInfo := GameErrorInformation{
				Message: "The game has ended!",
			}

			responseData, err := json.Marshal(MessageResponse[GameErrorInformation]{
				Type: "Error",
				Data: errorInfo,
			})

			if err != nil {
				s.Write([]byte("Error sending response"))
				return
			}

			s.Write(responseData)
			return
		}

		var data MessageFromPlayer

		err := json.Unmarshal(msg, &data)
		if err != nil {
			fmt.Println("Error getting message", err)
			return
		}

		ps := gs.GetPlayerSession(data.PlayerId)
		if ps == nil {
			s.Write([]byte("Cannot find player session"))
			return
		}

		newBoard := ps.Game.ValidateNewCell(data.Coors, data.Value)
		ps.Game.PrintBoard()
		resData := MessageResponse[[][]sudoku.Cell]{
			Type: "board",
			Data: newBoard,
		}
		newBoardData, err := json.Marshal(resData)
		if err != nil {
			fmt.Println("Error getting sudoku for player", ps.PlayerId, err)
			return
		}

		s.Write(newBoardData)

		psInfo := gs.GetPlayersProgressInformation()
		bcData := MessageResponse[[]OpponentPlayerInformation]{
			Type: "players",
			Data: psInfo,
		}
		playersInfo, err := json.Marshal(bcData)
		if err != nil {
			fmt.Println("Error broadcasting data to players", err)
			return
		}
		m.Broadcast(playersInfo)

		if ps.Game.IsValidBoard() {
			gs.onGoing = false
			gs.Winner = ps.PlayerName

			gInfo := GameEndedInformation{
				WinnerName:  gs.Winner,
				PlayersData: psInfo,
			}

			responseData, err := json.Marshal(MessageResponse[GameEndedInformation]{
				Type: "GameEnded",
				Data: gInfo,
			})

			if err != nil {
				fmt.Println("Error on parsing gameInfo", err)
				s.Write([]byte("Cannot get response"))
				return
			}

			m.Broadcast(responseData)
		}
	})

	fmt.Printf("Server running on port %s\r\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil); err != nil {
		fmt.Println("Error running server", err, PORT)
		panic("Cannot run server")
	}
}
