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
	// clean up private Key
	// send eth to the s.Client

	// nonce, err := ec.PendingNonceAt(context.TODO(), sender.PublicKey())
	// gasPrice, err := ec.SuggesGasPrice(context.TODO())
	// if err != nil {
	// return zero, err
	// }
	// s := types.NewEIP155Signer(chainId)
	// var amount big.Int
	// amount.SetInt64(amountToSend)
	// var gasLimit big.Int
	// gasLimit.SetInt64(1210000)
	// data := common.FromHex("0x")
	// t := types.NewTransaction(none, recepient, &amount, &gasLimit, gasPrice, data)
	// nt, err := types.SignTx(t,s, sender, GetKey())
	// if err != nil {
	//		return zero, err
	//	}
	// err := ec.SendTransaction(context.TODO(), nt)
	// return nt.Hash(), errrs
	return nil, nil
}
