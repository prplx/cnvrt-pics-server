package communicator

import (
	"strconv"

	"github.com/prplx/lighter.pics/internal/types"
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

func (c *Communicator) SendStartProcessing(jobID, fileID int, fileName string) error {
	return c.client.Trigger(channelName(jobID), "processing", types.AnyMap{
		"event":    "started",
		"fileName": fileName,
		"fileId":   fileID,
	})
}

func (c *Communicator) SendErrorProcessing(jobID, fileID int, fileName string) error {
	return c.client.Trigger(channelName(jobID), "processing", types.AnyMap{
		"event":    "error",
		"fileName": fileName,
		"fileId":   fileID,
	})
}

func (c *Communicator) SendSuccessProcessing(jobID int, result types.SuccessResult) error {
	return c.client.Trigger(channelName(jobID), "processing", types.AnyMap{
		"event":          "success",
		"fileId":         result.SourceFileID,
		"sourceFile":     result.SourceFileName,
		"targetFile":     result.TargetFileName,
		"sourceFileSize": strconv.FormatInt(result.SourceFileSize, 10),
		"targetFileSize": strconv.FormatInt(result.TargetFileSize, 10),
		"width":          result.Width,
		"height":         result.Height,
	})
}

func channelName(jobID int) string {
	return "cache-" + strconv.Itoa(jobID)
}
