package tcpserver

import (
	"math/big"

	. "github.com/wcrbrm/gethapi-example/server/database"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SendEthResponse struct {
	Status string              `json:"status"`
	Data   SendEthResponseBody `json:"data"`
}

type GetLastResponse struct {
	Status string                `json:"status"`
	Since  *big.Int              `json:"since"`
	Data   []GetLastResponseBody `json:"data"`
}
