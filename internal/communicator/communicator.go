package communicator

import (
	"strconv"

	"github.com/prplx/lighter.pics/internal/types"
	"github.com/pusher/pusher-http-go/v5"
)

type Communicator struct {
	client pusher.Client
}

func NewCommunicator(config *types.Config) *Communicator {
	client := pusher.Client{
		AppID:   config.Pusher.AppID,
		Key:     config.Pusher.Key,
		Secret:  config.Pusher.Secret,
		Cluster: config.Pusher.Cluster,
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

func (c *Communicator) SendStartArchiving(jobID int) error {
	return c.client.Trigger(channelName(jobID), "archiving", types.AnyMap{
		"event": "started",
	})
}

func (c *Communicator) SendErrorArchiving(jobID int) error {
	return c.client.Trigger(channelName(jobID), "archiving", types.AnyMap{
		"event": "error",
	})
}

func (c *Communicator) SendSuccessArchiving(jobID int, path string) error {
	return c.client.Trigger(channelName(jobID), "archiving", types.AnyMap{
		"event": "success",
		"path":  path,
	})
}

func channelName(jobID int) string {
	return "cache-" + strconv.Itoa(jobID)
}
