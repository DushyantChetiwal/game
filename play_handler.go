package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func enterMatchmaking(playerID string) string {
	_, err := rdb.LPush(context.Background(), "matchmaking_queue", playerID).Result()
	if err != nil {
		log.Printf("Failed to add player %s to matchmaking queue: %v", playerID, err)
		return ""
	}
	opponentID := ""
	// Wait for the match to be made in Redis
	for i := 0; i < 100; i++ {
		// Check if the player has been matched
		matched, err := rdb.Get(context.Background(), playerID).Result()
		if err != nil {
			log.Printf("Failed to check if player %s is matched: %v", playerID, err)
		} else {
			opponentID = matched
			break
		}

		// Sleep for a short time before checking again
		time.Sleep(500 * time.Millisecond)
	}

	return opponentID
}

func play(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Failed to create websocket: %v", err)
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("Error while closing connection %v", err)
		}
	}()

	playerID := uuid.New().String() + conn.RemoteAddr().String()
	// Add the player to the matchmaking queue in Redis

	opponentID := enterMatchmaking(playerID)
	if opponentID == "" {
		return
	}
	// Now that the match has been made, proceed with the game logic

	err = conn.WriteMessage(websocket.TextMessage, []byte("Please enter your choice rock, paper or scissor:"))
	if err != nil {
		log.Printf("Failed to get user input: %v", err)
		return
	}

	_, p1Choice, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read player response value: %v", err)
	}

	err = rdb.Set(context.Background(), "player_choice"+playerID, string(p1Choice), 10*time.Second).Err()
	if err != nil {
		log.Printf("Failed to set value: %v", err)
	}
	choice := ""
	for i := 0; i < 100; i++ {
		// Check if the opponent player has made the choice
		opponent_choice, err := rdb.Get(context.Background(), "player_choice"+opponentID).Result()
		if err != nil {
			log.Printf("Failed to check if opponent player %s has made a choice: %v", opponentID, err)
		} else {
			choice = opponent_choice
			break
		}
		// Sleep for a short time before checking again
		time.Sleep(500 * time.Millisecond)
	}

	gameLogic(conn, string(p1Choice), choice)
}
