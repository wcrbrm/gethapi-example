package tcpserver

import blockchain "github.com/wcrbrm/gethapi-example/server/blockchain"

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SendEthResponse struct {
	Status string                         `json:"status"`
	Data   blockchain.SendEthResponseBody `json:"data"`
}

type GetLastResponse struct {
	Status string                         `json:"status"`
	Data   blockchain.GetLastResponseBody `json:"data"`
}
