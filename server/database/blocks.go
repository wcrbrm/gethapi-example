package database

import (
	"log"
	"math/big"
)

func (s *DbClient) GetLastBlock() big.Int {
	var number big.Int
	err := s.DB.Get(&number, `SELECT COALESCE(lb.number, -1) AS number FROM view_last_block lb`)
	if err != nil {
		log.Fatal("Error retrieving last block in database. Please check view_last_block exists", err)
	}
	return number
}

// recursive function to accumulate hashed of blocks
// to increase confirmations on them
func (s *DbClient) GetBlocksToConfirm(parentHash string, depth int) []string {
	// if depth == 0 { return []string{} }
	return []string{}
}

func (s *DbClient) SaveBlock(number big.Int,
	hash string,
	parentHash string,
	properties map[string]interface{},
	transactions map[string]interface{}) {

	// get hashes of N previous blocks,
	// hashes := s.GetBlocksToConfirm(parentHash, s.nConfirmations)

	tx, err := s.DB.Begin()
	// 1) insert block with NamedExec?
	// 2) insert transactions with NamedExec
	// 3) insert or update account balances
	// 4) update number of confirmations - on previous blocks
	err = tx.Commit()
	if err != nil {
		log.Fatal("Initial Allocation Error", err)
	}
}
