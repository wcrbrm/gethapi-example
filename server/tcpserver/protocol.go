package tcpserver

import . "github.com/wcrbrm/gethapi-example/server/database"

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
	Data   []GetLastResponseBody `json:"data"`
}
