package tcpserver

import (
	"log"
)

type Callbacks struct {
	OnNewConnection        func(clientUid string)
	OnConnectionTerminated func(clientUid string)
}

func onNewConnection(clientUid string) {
	log.Printf("[server][%s] New connection", clientUid)
}

func onConnectionTerminated(clientUid string) {
	log.Printf("[server][%s] Terminated connection", clientUid)
}
