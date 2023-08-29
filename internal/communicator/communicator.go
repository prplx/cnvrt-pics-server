package communicator

import (
	"github.com/gorilla/websocket"
)

type Communicator struct {
	clients map[string]map[*websocket.Conn]bool
}

func NewCommunicator() *Communicator {
	return &Communicator{}
}
