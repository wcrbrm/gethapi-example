package database

import (
	"log"
	"math/big"
	"strings"
)

func (s *DbClient) GetLastBlock() big.Int {
	var number big.Int = *big.NewInt(-1)
	rows, err := s.DB.Query(`SELECT COALESCE(lb.number, -1) AS number FROM view_last_block lb`)
	if err != nil {
		log.Fatal("[blocks] Error retrieving last block in database. Please check view_last_block exists ", err)
	}
	if rows.Next() {
		err = rows.Scan(&number)
	}
	// log.Println("[blocks] GetLastBlock() returns", number.String())
	return number
}

// recursive function to accumulate hashed of blocks
// to increase confirmations on them
func (s *DbClient) GetBlocksToConfirm(parentHash string, depth int) []string {
	if depth == 0 || parentHash == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return []string{}
	}
	log.Println("Seeking for block with hash", parentHash)
	return []string{}
}

func FieldsFromMap(data map[string]string) (string, string) {
	fields := make([]string, 0)
	values := make([]string, 0)
	for key := range data {
		fields = append(fields, key)
		values = append(values, ":"+key)
	}
	return strings.Join(fields, ","), strings.Join(values, ",")
}

func (s *DbClient) SaveBlock(number *big.Int,
	parentHash string,
	blockProps map[string]string,
	transactions []map[string]string) {

	// get balances of accounts from this block
	// accounts := map[string]*big.Int{}
	// for _, t := range transactions {
	// }

	// get hashes of N previous blocks, to increase confirmations
	hashes := s.GetBlocksToConfirm(parentHash, s.nConfirmations)
	log.Println("[block-save] blocks to add confirmation: ", hashes)
	tx := s.DB.MustBegin()
	if len(hashes) > 0 {
		// 1) update number of confirmations - on previous blocks
		hashIndexes := make([]string, 0)
		hashMap := map[string]string{}
		for index, hash := range hashes {
			alias := "hash" + string(index)
			hashIndexes = append(hashIndexes, ":"+alias)
			hashMap[alias] = hash
		}
		_, errConfirmations := tx.NamedExec(
			" UPDATE blocks SET confirmations=confirmations+1 "+
				" WHERE hash IN ("+strings.Join(hashIndexes, ",")+")",
			hashMap)
		if errConfirmations != nil {
			log.Println("[block-save] Updating confirmations error: ", errConfirmations)
		}
	}

	// 2) insert block record
	fields, values := FieldsFromMap(blockProps)
	_, errBlock := tx.NamedExec(
		"INSERT INTO blocks ("+fields+") "+
			" VALUES ("+values+");", blockProps)
	if errBlock != nil {
		log.Println("[block-save] Inserting Block Error: ", errBlock)
	}

	// 3) insert transaction records, accumulate account balances
	for _, t := range transactions {
		// TODO:
		log.Println("[blocks] Transaction to be inserted", t)
	}

	// 4) insert or update account balances
	// TODO:

	err := tx.Commit()
	if err != nil {
		log.Fatal("[blocks] Initial Allocation Error", err)
	}
}
