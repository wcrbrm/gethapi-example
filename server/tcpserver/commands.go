package tcpserver

import (
	"encoding/json"
	"log"

	database "github.com/wcrbrm/gethapi-example/server/database"
)

func (c *Client) CommandGetLast() {
	data, err := c.chain.GetLast(*c.lastBlock)
	if err != nil {
		log.Println(err)
	}
	response := GetLastResponse{"success", c.lastBlock, *data}
	for _, row := range *data {
		// for each row, check if the block is
		if c.lastBlock.Cmp(&row.Number) < 0 {
			c.lastBlock.Set(&row.Number)
		}
	}
	json, _ := json.Marshal(response)
	c.conn.Write([]byte(string(json) + "\r\n"))
}

func (c *Client) CommandSendEth(command string) {
	subs := command[len("SendEth"):]
	// parse payload here
	var payload database.SendEthRequest
	errPayload := json.Unmarshal([]byte(subs), &payload)
	if errPayload != nil {
		log.Printf("[command][%s] SendEth payload parser error: %s", c.Uid, errPayload)
		c.SendResponse(Response{"error", errPayload.Error()})
	} else {
		log.Printf("[command][%s] SendEth payload: '%s'", c.Uid, payload)
		data, err := c.chain.SendEth(payload)
		if err != nil {
			log.Printf("[command][%s] '%s'", c.Uid, err)
			c.SendResponse(Response{"error", err.Error()})
		} else {
			response := SendEthResponse{"success", *data}
			json, _ := json.Marshal(response)
			c.conn.Write([]byte(string(json) + "\r\n"))
		}
	}
}
