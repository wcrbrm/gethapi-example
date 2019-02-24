package blockchain

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"time"
)

func (s *BlockchainClient) SyncBlock(num *big.Int) {
	log.Println("[blockchain] syncing block #", num.String())

	block, err := s.Client.BlockByNumber(context.Background(), num)
	if err != nil {
		log.Fatal("[blockchain] block error: ", num.String(), err)
	}

	h := block.Header()
	logsBloom, _ := h.Bloom.MarshalText()
	blockProps := map[string]string{
		"number":           h.Number.String(),
		"hash":             block.Hash().String(),
		"confirmations":    "0",
		"timestamp":        h.Time.String(),
		"parentHash":       h.ParentHash.String(), // common.Hash
		"nonce":            string(h.Nonce.Uint64()),
		"sha3Uncles":       h.UncleHash.String(),
		"logsBloom":        string(logsBloom),
		"transactionsRoot": h.TxHash.String(),
		"stateRoot":        h.Root.String(),
		"receiptsRoot":     h.ReceiptHash.String(),
		"miner":            h.Coinbase.Hex(),
		"difficulty":       h.Difficulty.String(),
		"extraData":        string(h.Extra),
		"gasLimit":         strconv.FormatUint(h.GasLimit, 10),
		"gasUsed":          strconv.FormatUint(h.GasUsed, 10),
		"mixhash":          h.MixDigest.String(),
	}
	txProps := []map[string]string{}

	// t := block.transactions
	s.DB.SaveBlock(num,
		h.ParentHash.String(),
		blockProps,
		txProps)
}

func (s *BlockchainClient) WatchBlocks() {
	log.Println("[blockchain] IN SYNC")

	d := 3 * time.Second
	for {
		time.Sleep(d)
		log.Println("[blockchain] watching blocks carefully")
		s.EnsureSynced(false)
	}
}

// on the first run we need to be verbose,
// on every second run - we don't neet to be verbose
func (s *BlockchainClient) EnsureSynced(firstRun bool) {
	verbose := firstRun
	dbBlockNum := s.DB.GetLastBlock()

	// 1. Get the last block in the database
	if verbose {
		log.Println("[blockchain] Last block in database:", dbBlockNum.String())
	}
	if firstRun && dbBlockNum.Cmp(big.NewInt(0)) <= 0 {
		// If there are no blocks, check genesis and create accounts funded with initial allocation
		s.DB.InitialAllocation(s.GetGenesisAllocation())
	}

	// 2. Get the last block in the blockchain
	header, err := s.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("[blockchain] getting latest block header error: ", err)
	}
	var chainNum big.Int = *header.Number
	if verbose {
		log.Println("[blockchain] block header # ", chainNum.String())
	}

	// from the last block in database (not included)
	// to the last block in the blockchain (included)
	// run importer of the block
	one := big.NewInt(1)
	startWith := dbBlockNum.Add(&dbBlockNum, one)
	for num := new(big.Int).Set(startWith); num.Cmp(&chainNum) <= 0; num.Add(num, one) {
		s.SyncBlock(num)
	}
}
