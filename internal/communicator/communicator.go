package communicator

import (
	"github.com/pusher/pusher-http-go/v5"
)

type Communicator struct {
	client pusher.Client
}

func NewCommunicator() *Communicator {
	client := pusher.Client{
		AppID:   "1528434",
		Key:     "238e350521ef2c91b881",
		Secret:  "ada76c58d004a6db5abe",
		Cluster: "eu",
		Secure:  true,
	}

	return &Communicator{
		client: client,
	}
}

func (c *Communicator) SendStartProcessing(jobID, fileName string) error {
	return c.client.Trigger(jobID, "processing", map[string]string{
		"event": "started",
		"file":  fileName,
	})
}

func (c *Communicator) SendErrorProcessing(jobID, fileName string) error {
	return c.client.Trigger(jobID, "processing", map[string]string{
		"event": "error",
		"file":  fileName,
	})
}
