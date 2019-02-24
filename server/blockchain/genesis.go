package blockchain

import (
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core"
)

func (s *BlockchainClient) ReadGenesis() *core.Genesis {
	genesis := new(core.Genesis)
	if len(s.GenesisPath) == 0 {
		return genesis // it is ok NOT to have a Genesis file provided
	}
	file, err := os.Open(s.GenesisPath)
	if err != nil {
		log.Fatal("[genesis] Failed to read genesis file", err)
	}
	defer file.Close()
	log.Println("[genesis] Found genesis file, reading initial distribution from ", s.GenesisPath)

	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		log.Fatal("[genesis] Invalid genesis file", err)
	}
	return genesis
}

func (s *BlockchainClient) GetGenesisAllocation() map[string]big.Int {
	genesis := s.ReadGenesis()
	m := map[string]big.Int{}
	for address, account := range genesis.Alloc {
		var balance big.Int = *account.Balance
		log.Println("[genesis] Initial balance: ", address.Hex(), balance.String())
		m[address.Hex()] = balance
	}
	return m
}
