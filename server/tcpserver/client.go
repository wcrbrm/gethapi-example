package tcpserver

import (
	"net"
	"bufio"
	"io"
	"log"
        "strings"
        "time"
)

type ConnectionEventType string

const (
	CONNECTION_EVENT_TYPE_NEW_CONNECTION        ConnectionEventType = "new_connection"
	CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED ConnectionEventType = "connection_terminated"
	CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR ConnectionEventType = "general_error"
)

// Client holds info about connection
type Client struct {
	Uid  string 
	conn net.Conn
	onConnectionEvent func(c *Client, eventType ConnectionEventType, e error) /* function for handling new connections */
	// onDataEvent func(c *Client, data []byte)
	onTextEvent func(c *Client, str string)
}


func NewClient(conn net.Conn, onConnectionEvent func(c *Client,eventType ConnectionEventType, e error), onTextEvent func(c *Client, str string)) *Client {
	return  &Client{
		conn: conn,
		onConnectionEvent: onConnectionEvent,
		// onDataEvent: onDataEvent,
	        onTextEvent: onTextEvent,
	}
}


// Read client data from channel
func (c *Client) listen() {
	timeout := 120 * time.Second
	err := c.conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(c.conn)
	// buf := make([]byte, 1024)
	for {
		str, err := reader.ReadString('\n')
		switch err {
		case io.EOF:
			// connection terminated
			c.conn.Close()
			c.onConnectionEvent(c,CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED, err)
			return
		case nil:
			// new data available
			c.onTextEvent(c, strings.TrimSpace(str))
		default:
			log.Fatalf("[server] Receive data failed:%s", err)
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

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	if c.conn!=nil {
		return c.conn.Close()
	}
	return nil
}
