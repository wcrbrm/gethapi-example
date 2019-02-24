package tcpserver

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net"
	"strings"
	_ "time"

	. "github.com/wcrbrm/gethapi-example/server/blockchain"
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
	lastBlock         big.Int
	chain             *BlockchainClient
	onConnectionEvent func(c *Client, eventType ConnectionEventType, e error) /* function for handling new connections */
}

func NewClient(conn net.Conn,
	chain *BlockchainClient,
	onConnectionEvent func(c *Client, eventType ConnectionEventType, e error)) *Client {
	return &Client{
		conn:              conn,
		chain:             chain,
		lastBlock:         *big.NewInt(-1),
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
				c.CommandSendEth(command)
			} else if command == "GetLast" {
				log.Printf("[server][%s] GetLast: '%s', since %s", c.Uid, command, c.lastBlock.String())
				c.CommandGetLast()
			} else if command != "" {
				log.Printf("[server][%s] Data received: '%s', ignored", c.Uid, command)
				c.SendResponse(Response{"error", "No Command"})
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
