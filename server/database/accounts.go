package database

import "log"

func (s *DbClient) InitialAllocation(accounts map[string]int64) {
	tx, err := s.DB.Begin()
	for accountId, amount := range accounts {
		arg := map[string]interface{}{
			"address": accountId,
			"balance": amount,
		}
		_, err = tx.Exec("INSERT INTO addresses (address, balance) "+
			"VALUES (:address, :balance)", arg)
		if err != nil {
			log.Fatal("Initial Allocation Address Error", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("Initial Allocation Error", err)
	}
}