package blockchain

import (
	"log"

	. "github.com/wcrbrm/gethapi-example/server/database"
)

func (s *BlockchainClient) GetLast(since int) *GetLastResponseBody {
	log.Println("[api] GetLast Since", since)
	return s.DB.GetLastTransactions(since)
}

func (s *BlockchainClient) SendEth(req *SendEthRequest) (*SendEthResponseBody, error) {
	log.Println("[api] SendEth")
	// check if address is valid
	// check account balance before sending. if nothing on this account, sorry, nothing can be sent
	// clean private Key
	return nil, nil
}
