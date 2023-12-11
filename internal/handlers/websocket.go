package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/contrib/websocket"
)

func (h *Handlers) HandleWebsocket(c *websocket.Conn) {
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
			}
			return
		} else {
			fmt.Println("websocket message received of type", messageType)
		}
	}
}
