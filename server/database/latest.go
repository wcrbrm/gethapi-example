package database

func (s *DbClient) GetLastTransactions(sinceBlock int) *GetLastResponseBody {
	// Последними считаются поступления,
	// которые либо еще не запрашивались данным методом,
	// либо имеют < 3 подтверждений на момент запроса.

	// - **date** - дата поступления.
	// - **address** - адрес, на который был произведен перевод
	// - **amount** - сумма перевода в ETH.
	// - **confirmations** - количество подтверждений транзакции. (Если
	// err := s.DB.Get(&number, `SELECT COALESCE(lb.number, -1) AS number FROM view_last_block lb`)

	// getting (blocks with < 3 confirmations) or (blocks)
	return nil
}
