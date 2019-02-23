package tcpserver

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	uuid "github.com/satori/go.uuid"
	blockchain "github.com/wcrbrm/gethapi-example/server/blockchain"
	database "github.com/wcrbrm/gethapi-example/server/database"
)

// TCP Server
type TcpServer struct {
	address     string
	connLock    sync.RWMutex
	connections map[string]*Client
	callbacks   Callbacks
	listener    net.Listener
	db          *database.DbClient
	chain       *blockchain.BlockchainClient
}

func (s *TcpServer) onConnectionEvent(c *Client, eventType ConnectionEventType, e error) {
	switch eventType {
	case CONNECTION_EVENT_TYPE_NEW_CONNECTION:
		s.connLock.Lock()
		u1, _ := uuid.NewV4()
		uidString := u1.String()
		c.Uid = uidString
		s.connections[uidString] = c
		s.connLock.Unlock()
		//log.Println(eventType , " ,  uid:", c.Uid, " , ip: ", c.Conn().RemoteAddr().String())
		if s.callbacks.OnNewConnection != nil {
			s.callbacks.OnNewConnection(uidString)
		}
	case CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED, CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR:
		//log.Println(eventType , " ,  uid:", c.Uid, " , ip: ", c.Conn().RemoteAddr().String(), " , error: ", e.Error())
		s.connLock.Lock()
		delete(s.connections, c.Uid)
		s.connLock.Unlock()
		if s.callbacks.OnConnectionTerminated != nil {
			s.callbacks.OnConnectionTerminated(c.Uid)
		}
	}
}

// Start network Server
func (s *TcpServer) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("[server] Error starting TCP Server: ", err)
	}
	for {
		conn, _ := s.listener.Accept()
		client := NewClient(conn, s.chain, s.onConnectionEvent)
		s.onConnectionEvent(client, CONNECTION_EVENT_TYPE_NEW_CONNECTION, nil)
		go client.listen()
	}
}

// Creates new tcp Server instance
func NewServer(address string, db *database.DbClient, chain *blockchain.BlockchainClient,
	callbacks Callbacks) *TcpServer {
	log.Println("[server] creating server with address", address)
	s := &TcpServer{
		address:   address,
		callbacks: callbacks,
		db:        db,
		chain:     chain,
	}
	s.connections = make(map[string]*Client)
	return s
}

func (s *TcpServer) SendDataByClientId(clientUid string, data []byte) error {
	if s.connections[clientUid] != nil {
		return s.connections[clientUid].Send(data)
	} else {
		return errors.New(fmt.Sprint("[server] no connection with uid ", clientUid))
	}
}

func (s *TcpServer) Close() {
	log.Println("[server] TcpServer.Close()")
	log.Println("[server] s.connections length: ", len(s.connections))
	for k := range s.connections {
		fmt.Printf("key[%s]\n", k)
		s.connections[k].Close()
	}
	s.listener.Close()
}

func NewTcpServer(db *database.DbClient, chain *blockchain.BlockchainClient) *TcpServer {
	HOST, ok := os.LookupEnv("HOST")
	if !ok {
		HOST = "127.0.0.1"
	}
	PORT, ok := os.LookupEnv("PORT")
	if !ok {
		PORT = "9091"
	}
	address := HOST + ":" + PORT
	return NewServer(address, db, chain, Callbacks{
		OnConnectionTerminated: onConnectionTerminated,
		OnNewConnection:        onNewConnection,
	})
}
