package blockchain

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func (s *BlockchainClient) SyncBlock(num *big.Int) {
	log.Println("[blockchain] syncing block #", num.String(), " chainID=", s.ChainID.String())

	block, err := s.Client.BlockByNumber(context.Background(), num)
	if err != nil {
		log.Fatal("[blockchain] block error: ", num.String(), err)
	}

	h := block.Header()
	log.Println("[blockchain] block #", num.String(), "time=", h.Time, "size=", len(block.Transactions()))

	logsBloom, _ := h.Bloom.MarshalText()
	blockProps := map[string]interface{}{
		"number":           h.Number.String(),
		"hash":             block.Hash().String(),
		"confirmations":    "0",
		"timestamp":        h.Time.String(),
		"parenthash":       h.ParentHash.String(), // common.Hash
		"nonce":            strconv.FormatUint(h.Nonce.Uint64(), 10),
		"sha3uncles":       h.UncleHash.String(),
		"logsbloom":        string(logsBloom),
		"transactionsroot": h.TxHash.String(),
		"stateroot":        h.Root.String(),
		"receiptsroot":     h.ReceiptHash.String(),
		"miner":            h.Coinbase.Hex(),
		"difficulty":       h.Difficulty.String(),
		"extradata":        "", // string(h.Extra),
		"gaslimit":         strconv.FormatUint(h.GasLimit, 10),
		"gasused":          strconv.FormatUint(h.GasUsed, 10),
		"mixhash":          h.MixDigest.String(),
	}

	txProps := []map[string]interface{}{}
	for txIndex, t := range block.Transactions() {
		vv, rr, ss := t.RawSignatureValues()
		var kvT = map[string]interface{}{
			"hash":             t.Hash().Hex(),
			"nonce":            strconv.FormatUint(t.Nonce(), 10),
			"blockhash":        block.Hash().String(),
			"blocknumber":      h.Number.String(),
			"transactionindex": txIndex,
			"to":               t.To().Hex(),
			"value":            t.Value().String(),
			"gas":              strconv.FormatUint(t.Gas(), 10),
			"gasprice":         t.GasPrice().String(),
			"v":                vv.String(),
			"r":                rr.String(),
			"s":                ss.String(),
		}
		if msg, err := t.AsMessage(types.NewEIP155Signer(s.ChainID)); err != nil {
			kvT["from"] = msg.From().Hex()
		}
		// log.Println("[blockchain] tx #", txIndex, kvT)
		txProps = append(txProps, kvT)
	}

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
		log.Println("[blockchain] waiting for more blocks")
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
	if firstRun && dbBlockNum.Cmp(big.NewInt(0)) < 0 {
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
