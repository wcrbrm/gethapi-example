package blockchain

import (
	"log"
	"time"
	_ "github.com/wcrbrm/gethapi-example/server/database"
)

func EnsureConnected() {
	log.Println("[blockchain] ensure connected")
}

func EnsureSynced() {
	log.Println("[blockchain] ensure synced")
}

func WatchBlocks() {
	d := 2*time.Second
	for range time.Tick(d) {
		log.Println("[blockchain] watching blocks carefully")
	}
}