package tcpserver

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strings"
	_ "time"

	blockchain "github.com/wcrbrm/gethapi-example/server/blockchain"
	. "github.com/wcrbrm/gethapi-example/server/database"
)

type ConnectionEventType string

const (
	CONNECTION_EVENT_TYPE_NEW_CONNECTION           ConnectionEventType = "new_connection"
	CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED    ConnectionEventType = "connection_terminated"
	CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR ConnectionEventType = "general_error"
)

// Client holds info about connection
type Client struct {
	Uid               string
	conn              net.Conn
	lastBlock         int
	chain             *blockchain.BlockchainClient
	onConnectionEvent func(c *Client, eventType ConnectionEventType, e error) /* function for handling new connections */
}

func NewClient(conn net.Conn,
	chain *blockchain.BlockchainClient,
	onConnectionEvent func(c *Client, eventType ConnectionEventType, e error)) *Client {
	return &Client{
		conn:              conn,
		chain:             chain,
		lastBlock:         -1,
		onConnectionEvent: onConnectionEvent,
	}
}

// Read client data from channel
func (c *Client) listen() {
	// timeout := 120 * time.Second
	// err := c.conn.SetReadDeadline(time.Now().Add(timeout))
	// if err != nil {
	//	panic(err)
	// }
	reader := bufio.NewReader(c.conn)
	for {
		str, err := reader.ReadString('\n')
		switch err {
		case io.EOF:
			// connection terminated
			c.conn.Close()
			c.onConnectionEvent(c, CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED, err)
			return
		case nil:
			// new data available
			command := strings.TrimSpace(str)
			if strings.HasPrefix(command, "SendEth") {
				log.Printf("[server][%s] SendEth received '%s'", c.Uid, command)
				subs := command[len("SendEth"):]
				// parse payload here
				var payload SendEthRequest
				errPayload := json.Unmarshal([]byte(subs), &payload)
				if errPayload != nil {
					log.Printf("[server][%s] SendEth payload parser error: %s", c.Uid, errPayload)
					c.SendResponse(Response{"error", errPayload.Error()})
				} else {
					log.Printf("[server][%s] SendEth payload: '%s'", c.Uid, command)
					// TODO: send client response from c.chain.SendEth(payload)
					c.SendResponse(Response{"success", "OK"})
				}

			} else if command == "GetLast" {
				log.Printf("[server][%s] GetLast: '%s', since %d", c.Uid, command, c.lastBlock)
				// TODO: send client response from c.chain.GetLast()
				c.SendResponse(Response{"success", "OK"})
			} else {
				log.Printf("[server][%s] Data received: '%s', ignored", c.Uid, command)
				c.SendResponse(Response{"error", "Nil"})
			}
		default:
			log.Fatalf("[server][%s] Receive data failed:%s", c.Uid, err)
			c.conn.Close()
			c.onConnectionEvent(c, CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR, err)
			return
		}
	}
}

// Send text message to client
func (c *Client) Send(message []byte) error {
	_, err := c.conn.Write(message)
	return err
}

// Send text message to client
func (c *Client) SendResponse(resp Response) error {
	json, _ := json.Marshal(resp)
	_, err := c.conn.Write([]byte(string(json) + "\r\n"))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
