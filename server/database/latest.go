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
	sql := "SELECT b.number, b.timestamp, b.confirmations, " +
		" t.hash hash, t.value, t.from, t.to " +
		" FROM blocks b, transactions t " +
		" WHERE b.number = t.blockNumber " +
		" AND (b.number < :confirmations OR b.number > :since)"

	rows, err := s.DB.NamedQuery(sql, arg)
	if err != nil {
		log.Fatal("[database] Error retrieving last transactions ", err)
	}

	result := new([]GetLastResponseBody)
	for rows.Next() {
		row := rows.Scan()
		log.Println("row: ", row)
		// - **date** - дата поступления.
		// - **address** - адрес, на который был произведен перевод
		// - **amount** - сумма перевода в ETH.
	}
	return result, nil
}
