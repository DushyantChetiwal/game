package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var rdb *redis.Client

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

func gameLogic(conn *websocket.Conn, p1Choice, p2Choice string) {
	if p2Choice == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("You won by abandonment"))
	}
	if p1Choice == p2Choice {
		conn.WriteMessage(websocket.TextMessage, []byte("You won"))
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte("You lost"))
	}
}

func play(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Failed to create websocket: %v", err)
		return
	}

	playerID := uuid.New().String() + conn.RemoteAddr().String()
	// Add the player to the matchmaking queue in Redis

	opponentID := enterMatchmaking(playerID)

	// Now that the match has been made, proceed with the game logic
	// ...
	err = conn.WriteMessage(websocket.TextMessage, []byte("Please enter your choice (0 or 1):"))
	if err != nil {
		log.Printf("Failed to get user input: %v", err)
		return
	}

	_, p1Choice, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read player response value: %v", err)
	}

	err = rdb.Set(context.Background(), "player_choice"+playerID, string(p1Choice), 0).Err()
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
	conn.Close()
}

func gameMaster() {
	for {
		queueLen, err := rdb.LLen(context.Background(), "matchmaking_queue").Result()
		if err != nil {
			log.Printf("Failed to get matchmaking queue length: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// If there are at least two players in the queue, match them
		if queueLen >= 2 {
			// Pop two players from the matchmaking queue
			player1, err := rdb.LPop(context.Background(), "matchmaking_queue").Result()
			if err != nil {
				log.Printf("Failed to pop player from matchmaking queue: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			player2, err := rdb.LPop(context.Background(), "matchmaking_queue").Result()
			if err != nil {
				log.Printf("Failed to pop player from matchmaking queue: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Add the matched players to the "matched_players" set
			_, err = rdb.Set(context.Background(), player1, player2, 0).Result()
			if err != nil {
				log.Printf("Failed to add players to matched players set: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			_, err = rdb.Set(context.Background(), player2, player1, 0).Result()
			if err != nil {
				log.Printf("Failed to add players to matched players set: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("match %v %v", player1, player2)
		}

		// Wait for a short time before checking again
		time.Sleep(1 * time.Second)
	}
}
func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()
	go gameMaster()
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port unspecified")
	}
	fmt.Println(portString)
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/play", play)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on:%v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
