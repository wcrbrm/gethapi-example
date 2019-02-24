package database

import (
	"log"
)

// recursive function to accumulate hashed of blocks
// to increase confirmations on them
func (s *DbClient) GetBlocksToConfirm(parentHash string, depth int) []string {
	if depth == 0 || parentHash == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return []string{}
	}
	// log.Println("[confirmations] Seeking for block with hash=", parentHash, "depth=", depth)
	sql := `SELECT b.number as n, b.parentHash as phash  FROM blocks b WHERE b.hash = :hash`
	rows, err := s.DB.NamedQuery(sql, map[string]interface{}{
		"hash": parentHash,
	})
	defer rows.Close()

	var res = []string{}
	if err != nil {
		log.Println("[confirmations] Error retrieving last block in database. Please check view_last_block exists ", err)
	} else if rows.Next() {
		var n int64
		var phash string
		err1 := rows.Scan(&n, &phash)
		if err1 != nil {
			log.Println("[confirmations] WARNING reading block error ", err1)
		} else {
			res = append(res, parentHash) // we are sure to increase confirmations on this block
			for _, blockHash := range s.GetBlocksToConfirm(phash, depth-1) {
				res = append(res, blockHash)
			}
		}
	}
	return res
}
