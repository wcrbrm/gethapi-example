package blockchain

import "log"

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
	Date          string `json:"date"`
	Address       string `json:"address"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmation"`
}

func (s *BlockchainClient) GetLast() *GetLastResponseBody {
	log.Println("[api] Getting Last")
	return nil
}

func (s *BlockchainClient) SendEth(req *SendEthRequest) *SendEthResponseBody {
	log.Println("[api] SendEth")
	return nil
}
