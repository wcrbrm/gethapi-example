package blockchain

import (
	"log"
	"os"
	"strconv"

	ethclient "github.com/ethereum/go-ethereum/ethclient"
	database "github.com/wcrbrm/gethapi-example/server/database"
)

type BlockchainClient struct {
	GethAddress          string
	NumConfirmations     int // # of confirmations for blocks to be kept forever
	NumLastConfirmations int // # of confirmations for "GetLast" output
	DB                   *database.DbClient
}

func NewGethConnection(db *database.DbClient) *BlockchainClient {
	GETH, ok := os.LookupEnv("GETH")
	if !ok {
		GETH = "http://localhost:8545"
	}
	log.Printf("[blockchain] connecting to %s\n", GETH)

	_, err := ethclient.Dial(GETH)
	if err != nil {
		log.Fatal("[blockchain] Couldn't connect to the GETH RPC server")
	}
	minConfirmations, ok := os.LookupEnv("GETH_MIN_CONFIRMATIONS")
	if !ok {
		minConfirmations = "6"
	}
	nConfirmations, _ := strconv.Atoi(minConfirmations)

	lastConfirmations, ok := os.LookupEnv("GETLAST_CONFIRMATIONS")
	if !ok {
		lastConfirmations = "3"
	}
	nLastConfirmations, _ := strconv.Atoi(lastConfirmations)
	return &BlockchainClient{
		GETH,
		nConfirmations,
		nLastConfirmations,
		db,
	}
}
