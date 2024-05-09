package main

import (
	"context"
	"log"
	"time"
)

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
			_, err = rdb.Set(context.Background(), player1, player2, 10*time.Second).Result()
			if err != nil {
				log.Printf("Failed to add players to matched players set: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			_, err = rdb.Set(context.Background(), player2, player1, 10*time.Second).Result()
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
