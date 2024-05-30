package main

import (
	"github.com/gorilla/websocket"
)

var choices = map[string]string{
	"rock":    "scissor",
	"paper":   "rock",
	"scissor": "paper",
}

func gameLogic(conn *websocket.Conn, p1Choice, p2Choice string) {
	if p2Choice == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("You won by abandonment"))
		return
	}
	winningCase, ok := choices[p1Choice]
	if ok {
		if p1Choice == p2Choice {
			conn.WriteMessage(websocket.TextMessage, []byte("It's a tie!"))
			return
		}

		// Check if player 1 wins
		if winningCase == p2Choice {
			conn.WriteMessage(websocket.TextMessage, []byte("You won"))
			return
		}
	}

	// Player 2 wins
	conn.WriteMessage(websocket.TextMessage, []byte("You lost"))
}
