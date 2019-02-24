package database

import (
	"log"
	"math/big"
)

func (s *DbClient) GetLastTransactions(sinceBlock big.Int) (*[]GetLastResponseBody, error) {
	arg := map[string]interface{}{
		"confirmations": s.nConfirmations,
		"since":         sinceBlock.String(),
	}
	sql := "SELECT b.timestamp, b.confirmations, " +
		" t.value * 1e-18 amount, t.to address, b.number" +
		" FROM blocks b, transactions t " +
		" WHERE b.number = t.blockNumber " +
		" AND (b.number < :confirmations OR b.number > :since)" +
		" ORDER BY b.number"

	rows, err := s.DB.NamedQuery(sql, arg)
	if err != nil {
		log.Fatal("[database] Error retrieving last transactions ", err)
	}

	result := []GetLastResponseBody{}
	for rows.Next() {
		var timestamp string
		var amount string
		var confirmations int
		var address string
		var number int64
		err := rows.Scan(&timestamp, &confirmations, &amount, &address, &number)
		result = append(result, GetLastResponseBody{
			Date:          timestamp,
			Address:       address,
			Amount:        amount,
			Confirmations: confirmations,
			Number:        *big.NewInt(number),
		})
		if err != nil {
			log.Println("Get Last Transactions Error: ", err)
		}
	}
	return &result, nil
}
