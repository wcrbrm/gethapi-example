package tcpserver

import (
	"log"
)

type Callbacks struct {
	OnNewConnection        func(clientUid string)
	OnConnectionTerminated func(clientUid string)
}

func onNewConnection(clientUid string) {
	log.Printf("[server][%s] new connection", clientUid)
}

func onConnectionTerminated(clientUid string) {
	log.Printf("[server][%s] terminated connection", clientUid)
}
