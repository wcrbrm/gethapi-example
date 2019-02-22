package tcpserver

import (
	"log"
)

type Callbacks struct {
	OnNewConnection func (clientUid string)
	OnConnectionTerminated func (clientUid string)
}

func onNewConnection(clientUid string) {
	log.Printf("[server] new connection %s", clientUid)
}

func onConnectionTerminated(clientUid string) {
	log.Printf("[server] terminated connection %s", clientUid)
}

