package blockchain

import (
	"context"
	"log"
	"time"
)

/*
[
  'number', 'hash', 'parentHash', 'nonce', 'sha3Uncles', 'logsBloom', 'transactionsRoot',
  'stateRoot', 'receiptsRoot', 'miner', 'difficulty', 'totalDifficulty', 'size',
  'extraData', 'gasLimit', 'gasUsed', 'timestamp', 'blockHash', 'blockNumber',
  'transactionIndex', 'from', 'to', 'value',  'gas', 'gasPrice', 'input',
  'mixhash', 'v', 'r', 's'
];
*/

func (s *BlockchainClient) WatchBlocks() {
	d := 2 * time.Second
	for range time.Tick(d) {
		log.Println("[blockchain] watching blocks carefully")
	}
}

func (s *BlockchainClient) EnsureSynced() {
	dbBlockNum := s.DB.GetLastBlock()

	// 1. Get the last block in the database
	log.Println("[blockchain] Last block in database:", dbBlockNum)
	if dbBlockNum <= 0 {
		// If there are no blocks, check genesis and create accounts funded with initial allocation
		// s.DB.InitialAllocation(s.GenesisPath)
	}

	// 2. Get the last block in the blockchain
	header, err := s.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("[blockchain] block header error", err)
	}
	log.Println("[blockchain] block header number ", header.Number)

	// from the last block in database (not included)
	// to the last block in the blockchain (included)
	// run importer of the block
	for num := dbBlockNum; num <= header.Number; num++ {
	}
}
