package database

import "math/big"

type SendEthRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Key    string `json:"key"`
}

type SendEthResponseBody struct {
	Address string `json:"address"`
	Tx      string `json:"tx"`
	Nonce   int    `json:"nonce"`
}

type GetLastResponseBody struct {
	Date          string  `json:"date"`
	Address       string  `json:"address"`
	Amount        string  `json:"amount"`
	Confirmations int     `json:"confirmation"`
	Number        big.Int `json:"number"`
}
