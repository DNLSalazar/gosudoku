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

type Message struct {
	Coors    sudoku.Coor `json:"coors"`
	Value    int         `json:"value"`
	PlayerId string      `json:"playerId"`
}

type PlayerSession struct {
	PlayerId string
	Game     sudoku.Sudoku
}

func (p *PlayerSession) ToPlayerSessionResponse() PlayerGameStartResponse {
	return PlayerGameStartResponse{
		PlayerId: p.PlayerId,
		Board:    p.Game.GetBoard(),
	}
}

type PlayerGameStartResponse struct {
	PlayerId string          `json:"playerId"`
	Board    [][]sudoku.Cell `json:"board"`
}

type GameSession struct {
	InitialBoard [][]sudoku.Cell
	Player1      *PlayerSession
	Player2      *PlayerSession
	Winner       string
	onGoing      bool
}

func (gs *GameSession) SlotAvailable() bool {
	return gs.Player1 == nil || gs.Player2 == nil
}

func (gs *GameSession) GetPlayerSession(id string) *PlayerSession {
	if gs.Player1 != nil && gs.Player1.PlayerId == id {
		return gs.Player1
	}
	if gs.Player2 != nil && gs.Player2.PlayerId == id {
		return gs.Player2
	}
	return nil
}

func (gs *GameSession) AddPLayerSession(ps *PlayerSession) error {
	if !gs.SlotAvailable() {
		return errors.New("No slots for player session")
	}

	if gs.Player1 == nil {
		gs.Player1 = ps
		return nil
	}

	if gs.Player2 == nil {
		gs.Player2 = ps
		return nil
	}

	return nil
}

func Server() {
	m := melody.New()
	game := sudoku.CreateNewSudoku(17)
	gs := GameSession{
		InitialBoard: game.GetBoard(),
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

		err := gs.AddPLayerSession(&pSession)
		if err != nil {
			fmt.Println("Error adding player session")
			w.Write([]byte("Cannot add player session to game"))
			return
		}

		data, err := json.Marshal(pSession.ToPlayerSessionResponse())
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
		var data Message

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

		fmt.Println("GAME SESSION FOUND", ps.PlayerId)

		newBoard := ps.Game.ValidateNewCell(data.Coors, data.Value)
		ps.Game.PrintBoard()
		newBoardData, err := json.Marshal(newBoard)
		if err != nil {
			fmt.Println("Error getting sudoku")
			return
		}

		s.Write(newBoardData)
		// m.Broadcast(newBoardData)
	})

	http.ListenAndServe(":5000", nil)
}
