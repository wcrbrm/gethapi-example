package database

import (
	"log"
	"math/big"
	"strconv"
	"strings"
)

func (s *DbClient) GetLastBlock() big.Int {
	var number big.Int = *big.NewInt(-1)
	rows, err := s.DB.Query(`SELECT COALESCE(lb.number, -1) AS number FROM view_last_block lb`)
	if err != nil {
		log.Fatal("[blocks] Error retrieving last block in database. Please check view_last_block exists ", err)
	}
	defer rows.Close()
	if rows.Next() {
		var num int64
		err1 := rows.Scan(&num)
		if err1 != nil {
			log.Fatal("[blocks] Reading last block", err1)
		}
		number = *big.NewInt(num)
	}
	// log.Println("[blocks] GetLastBlock() returns", number.String())
	return number
}

func FieldsFromMap(data map[string]interface{}) (string, string) {
	fields := make([]string, 0)
	values := make([]string, 0)
	for key := range data {
		fields = append(fields, "\""+key+"\"")
		values = append(values, ":"+key)
	}
	return strings.Join(fields, ","), strings.Join(values, ",")
}

func (s *DbClient) SaveBlock(number *big.Int,
	parentHash string,
	blockProps map[string]interface{},
	transactions []map[string]interface{}) {

	tx := s.DB.MustBegin()
	// create addresses of destinations that are not in the database yet
	for _, tProps := range transactions {
		addr := tProps["to"]
		arg := map[string]interface{}{
			"address": addr,
			"balance": "0",
		}
		// log.Println("[address] adding address", addr)
		// the error below is ignored very intentionally
		tx.NamedExec("INSERT INTO addresses (address, value) "+
			"VALUES (:address, :balance) ON CONFLICT DO NOTHING", arg)
	}

	// get hashes of N previous blocks, to increase confirmations
	hashes := s.GetBlocksToConfirm(parentHash, s.nConfirmations)
	if len(hashes) > 0 {
		// log.Println("[block-save] blocks to add confirmation: ", hashes)
		// 1) update number of confirmations - on previous blocks
		hashIndexes := make([]string, 0)
		hashMap := map[string]interface{}{}
		for index, hash := range hashes {
			alias := "hash" + strconv.Itoa(index)
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
	sql := "INSERT INTO blocks (" + fields + ") VALUES (" + values + ");"
	// log.Println("[block-save] block SQL: ", sql)
	// log.Println("[block-save] block Properties: ", blockProps)
	_, errBlock := tx.NamedExec(sql, blockProps)
	if errBlock != nil {
		log.Println("[block-save] Inserting Block Error: ", errBlock)
	}

	// 3) insert transaction records, accumulate account balances
	for _, tProps := range transactions {
		fieldsT, valuesT := FieldsFromMap(tProps)
		sql := "INSERT INTO transactions (" + fieldsT + ") VALUES (" + valuesT + ")"

		// log.Println("[block-save] transaction SQL: ", sql)
		// log.Println("[block-save] transaction Properties: ", tProps)
		_, errT := tx.NamedExec(sql, tProps)
		if errT != nil {
			log.Println("[transaction-save] Inserting Transaction Error: ", errT)
		}

		// the following part is probably completely wrong
		// as we need to wait for confirmations, and should not update it
		// on transaction discovery

		// 4a) deduct balance from sender
		var argFrom = map[string]interface{}{
			"increment": tProps["value"],
			"address":   tProps["from"],
		}
		tx.NamedExec("UPDATE addresses "+
			"SET value = value - :increment, nonce = nonce + 1 "+
			"WHERE address = :address", argFrom)

		// 4b) add values to receiver
		var argTo = map[string]interface{}{
			"increment": tProps["value"],
			"address":   tProps["to"],
		}
		tx.NamedExec("UPDATE addresses "+
			"SET value = value + :increment "+
			"WHERE address = :address", argTo)
	}

	err := tx.Commit()
	if err != nil {
		log.Fatal("[blocks] Initial Allocation Error", err)
	}
}
