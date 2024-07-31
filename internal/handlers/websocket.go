package handlers

import (
	"strconv"
	"time"

	"github.com/prplx/cnvrt/internal/types"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (h *Handlers) HandleWebsocket(c types.WebsocketConnection) {
	jobID := c.Params("jobID")
	jobIDInt, err := strconv.ParseInt(jobID, 10, 64)
	send := make(chan []byte, 8)

	if err != nil {
		c.Close()
		return
	}

	defer func() {
		h.services.Communicator.RemoveClient(jobIDInt)
	}()

	h.services.Communicator.AddClient(jobIDInt, c)

	go handleKeepAlive(c, send)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}

		send <- message
	}
}

func handleKeepAlive(c types.WebsocketConnection, messageChannel chan []byte) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Close()
	}()

	c.SetReadDeadline(time.Now().Add(pongWait))

	for {
		select {
		case <-messageChannel:
			c.SetReadDeadline(time.Now().Add(pongWait))
		case <-ticker.C:
			pingEvent := types.AnyMap{"event": "ping", "operation": "keepalive"}
			if err := c.WriteJSON(pingEvent); err != nil {
				return
			}
		}
	}
}
