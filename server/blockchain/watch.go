package blockchain

import (
	"context"
	"log"
	"math/big"
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
	for {
		time.Sleep(d)
		log.Println("[blockchain] watching blocks carefully")
		// s.EnsureSynced(false)
	}
}

// on the first run we need to be verbose,
// on every second run - we don't neet to be verbose
func (s *BlockchainClient) EnsureSynced(verbose bool) {
	dbBlockNum := s.DB.GetLastBlock()

	// 1. Get the last block in the database
	if verbose {
		log.Println("[blockchain] Last block in database:", dbBlockNum)
	}
	if dbBlockNum.Cmp(big.NewInt(0)) <= 0 {
		// If there are no blocks, check genesis and create accounts funded with initial allocation
		s.DB.InitialAllocation(s.GetGenesisAllocation())
	}

	// 2. Get the last block in the blockchain
	header, err := s.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("[blockchain] block header error", err)
	}
	var chainNum big.Int = *header.Number
	if verbose {
		log.Println("[blockchain] block header number ", chainNum)
	}

	// from the last block in database (not included)
	// to the last block in the blockchain (included)
	// run importer of the block
	var one = big.NewInt(1)
	for num := new(big.Int).Set(&dbBlockNum); num.Cmp(&chainNum) <= 0; num.Add(num, one) {
		log.Println("[blockchain] syncing block#", num)
	}
}
