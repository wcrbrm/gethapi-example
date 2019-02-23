package blockchain

import (
	"log"
	"time"
)

func (s *BlockchainClient) WatchBlocks() {
	d := 2 * time.Second
	for range time.Tick(d) {
		log.Println("[blockchain] watching blocks carefully")
	}
}

func (s *BlockchainClient) EnsureSynced() {
	log.Println("[blockchain] ensure synced")
}
