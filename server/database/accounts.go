package database

import (
	"log"
	"math/big"
)

func (s *DbClient) InitialAllocation(accounts map[string]big.Int) {

	log.Println("[database] INITIALIZING")

	_, errCleanup := s.DB.Exec("DELETE FROM addresses")
	if errCleanup != nil {
		log.Fatal("[database] Accounts clean up error: ", errCleanup)
	}

	tx := s.DB.MustBegin()
	for addr, amount := range accounts {
		arg := map[string]interface{}{
			"address": addr,
			"balance": amount.String(),
		}
		_, errTx := tx.NamedExec("INSERT INTO addresses (address, value) "+
			"VALUES (:address, :balance);", arg)
		if errTx != nil {
			log.Fatal("[database] Account insert error: ", errTx)
		}
	}
	err := tx.Commit()
	if err != nil {
		log.Fatal("[database] Initial Allocation Error: ", err)
	}
}
