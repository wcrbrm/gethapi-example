package blockchain

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/wcrbrm/gethapi-example/server/database"
)

func (s *BlockchainClient) GetLast(since big.Int) (*[]GetLastResponseBody, error) {
	log.Println("[api] GetLast since=", since.String(), " last confirmations=", s.NumLastConfirmations)
	return s.DB.GetLastTransactions(since, s.NumLastConfirmations)
}

func (s *BlockchainClient) SendEth(req SendEthRequest) (*SendEthResponseBody, error) {
	log.Println("[api] SendEth")
	// check if address is valid
	// check account balance before sending. if nothing on this account, sorry, nothing can be sent
	ec := s.Client

	// verify private Key
	prv, err := crypto.HexToECDSA(req.Key)
	if err != nil {
		log.Println("[send-eth] private key error", err)
		return nil, err
	}

	sender := common.HexToAddress(req.From)
	receiver := common.HexToAddress(req.To)
	log.Println("[send-eth] sender=", sender.Hex(), " receiver=", receiver.Hex())

	nonce, err := ec.PendingNonceAt(context.Background(), sender)
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println("[send-eth] gas price suggestion error", err)
		return nil, err
	}
	log.Println("[send-eth] nonce=", nonce, "gasPrice=", gasPrice.String(), "chainID=", s.ChainID)

	signer := types.NewEIP155Signer(s.ChainID)
	var amount big.Int
	amount.SetString(req.Amount, 10)
	var gasLimit uint64 = 1210000
	log.Println("[send-eth] amount=", amount.String(), "gasLimit=", gasLimit)

	// data := common.FromHex("0x")
	t := types.NewTransaction(nonce, receiver, &amount, gasLimit, gasPrice, nil)
	nt, errSign := types.SignTx(t, signer, prv)
	if errSign != nil {
		log.Println("[send-eth] failed to sign: ", errSign)
		return nil, errSign
	}
	errTx := ec.SendTransaction(context.Background(), nt)
	if errTx == nil {
		log.Println("[send-eth] transaction was sent, hash=", nt.Hash().Hex())
	} else {
		log.Println("[send-eth] failed to send: ", errTx)
	}
	return &SendEthResponseBody{Tx: nt.Hash().Hex()}, errTx
}
