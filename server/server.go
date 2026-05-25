package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DNLSalazar/gosudoku/sudoku"
	"github.com/olahol/melody"
)

const PORT = "1289"

func Server(initialCells int) {
	m := melody.New()
	game := sudoku.CreateNewSudoku(initialCells)
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
