package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/prplx/cnvrt/internal/types"
)

func (h *Handlers) HandleWebsocket(c types.WebsocketConnection) {
	jobID := c.Params("jobID")
	jobIDInt, err := strconv.Atoi(jobID)
	if err != nil {
		c.Close()
		return
	}

	defer func() {
		h.services.Communicator.RemoveClient(jobIDInt)
		c.Close()
	}()

	h.services.Communicator.AddClient(jobIDInt, c)

	for {
		messageType, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("read error:", err)
				h.services.Logger.PrintError(err, types.AnyMap{
					"message":     "error reading message",
					"messageType": messageType,
				})
			}
			return
		}
	}
}
