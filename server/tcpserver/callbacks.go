package tcpserver

import "log"

type Callbacks struct {
	OnNewConnection func (clientUid string)
	OnConnectionTerminated func (clientUid string)
	// OnDataReceived func (clientUid string, data []byte)
	OnTextReceived func (clientUid string, str string)
}

func onNewConnection(clientUid string) {
	log.Printf("[server] new connection %s", clientUid)
}

func onConnectionTerminated(clientUid string) {
	log.Printf("[server] terminated connection %s", clientUid)
}

func onTextReceived(clientUid string, str string) {
	log.Printf("[server] data received from %s: '%s'", clientUid, str)
}
