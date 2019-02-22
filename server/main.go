package main

import (
	tcpserver  "github.com/wcrbrm/gethapi-example/server/tcpserver"
	database   "github.com/wcrbrm/gethapi-example/server/database"
	blockchain "github.com/wcrbrm/gethapi-example/server/blockchain"
)

func main() {
	// Validate that we can access database (halt program execution otherwise)	
	database.EnsureConnected()
	// Validate that we can access GETH (halt program execution otherwise)	
	blockchain.EnsureConnected()
	// Make sure that blockchain is in sync with the database before moving forward
	blockchain.EnsureSynced()
	// Start thread of blocks watching
	go blockchain.WatchBlocks()
	// Start our main TCP server - to listen for the messages from the clients
	tcpserver.Listen()
}