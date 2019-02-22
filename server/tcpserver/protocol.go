package tcpserver

type SendEthRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Key    string `json:"key"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SEND ETH interfaces
type SendEthResponseBody struct {
	Address string `json:"address"`
	Tx      string `json:"tx"`
	Nonce   string `json:"nonce"`
}

type SendEthResponse struct {
	Status string              `json:"status"`
	Data   SendEthResponseBody `json:"data"`
}

// GET LAST interfaces
type GetLastResponseBody struct {
	Date          string `json:"date"`
	Address       string `json:"address"`
	Amount        int    `json:"amount"`
	Confirmations int    `json:"confirmation"`
}

type GetLastResponse struct {
	Status string              `json:"status"`
	Data   GetLastResponseBody `json:"data"`
}
