package handlers

import (
	"strconv"

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
		_, _, err := c.ReadMessage()
		if err != nil {
			return
		}
	}
}
