package main

import (
	blockchain "github.com/wcrbrm/gethapi-example/server/blockchain"
	database "github.com/wcrbrm/gethapi-example/server/database"
	tcpserver "github.com/wcrbrm/gethapi-example/server/tcpserver"
)

func main() {
	// Validate that we can access database (halt program execution otherwise)
	db := database.NewDatabaseClient()

	// Validate that we can access GETH (halt program execution otherwise)
	chain := blockchain.NewGethConnection(db)
	// Make sure that blockchain is in sync with the database before moving forward
	chain.EnsureSynced(true)

	// Start thread of blocks watching
	go chain.WatchBlocks()
	// Start our main TCP server - to listen for the messages from the clients
	tcpserver.NewTcpServer(db, chain).Start()
}
