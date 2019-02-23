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
	NumLastConfirmations int // # of confirmations for "GetLast" output
	DB                   *database.DbClient
	Client               *ethclient.Client
	GenesisPath          string
}

func NewGethConnection(db *database.DbClient) *BlockchainClient {
	GETH, ok := os.LookupEnv("GETH")
	if !ok {
		GETH = "http://localhost:8545"
	}
	log.Printf("[blockchain] connecting to %s\n", GETH)

	client, err := ethclient.Dial(GETH)
	if err != nil {
		log.Fatal("[blockchain] Couldn't connect to the GETH RPC server")
	}
	lastConfirmations, ok := os.LookupEnv("GETLAST_CONFIRMATIONS")
	if !ok {
		lastConfirmations = "3"
	}
	nLastConfirmations, _ := strconv.Atoi(lastConfirmations)

	genesisPath, ok := os.LookupEnv("GETH_GENESIS_PATH")
	if !ok {
		genesisPath = ""
	}
	return &BlockchainClient{
		GETH,
		nLastConfirmations,
		db,
		client,
		genesisPath,
	}
}
