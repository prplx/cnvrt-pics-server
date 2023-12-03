package communicatorpusher

import (
	"fmt"
	"strconv"

	"github.com/prplx/lighter.pics/internal/types"
	"github.com/pusher/pusher-http-go/v5"
)

type Communicator struct {
	client pusher.Client
}

const (
	ProcessingEvent = "processing"
	ArchivingEvent  = "archiving"
	FlushedEvent    = "flushed"
	StartedEvent    = "started"
	ErrorEvent      = "error"
	SuccessEvent    = "success"
)

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
	return c.client.Trigger(channelName(jobID), ProcessingEvent, types.AnyMap{
		"event":    StartedEvent,
		"fileName": fileName,
		"fileId":   fileID,
	})
}

func (c *Communicator) SendErrorProcessing(jobID, fileID int, fileName string) error {
	return c.client.Trigger(channelName(jobID), ProcessingEvent, types.AnyMap{
		"event":    ErrorEvent,
		"fileName": fileName,
		"fileId":   fileID,
	})
}

func (c *Communicator) SendSuccessProcessing(jobID int, result types.SuccessResult) error {
	return c.client.Trigger(channelName(jobID), ProcessingEvent, types.AnyMap{
		"event":          SuccessEvent,
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
	return c.client.Trigger(channelName(jobID), ArchivingEvent, types.AnyMap{
		"event": StartedEvent,
	})
}

func (c *Communicator) SendErrorArchiving(jobID int) error {
	return c.client.Trigger(channelName(jobID), ArchivingEvent, types.AnyMap{
		"event": ErrorEvent,
	})
}

func (c *Communicator) SendSuccessArchiving(jobID int, path string) error {
	return c.client.Trigger(channelName(jobID), ArchivingEvent, types.AnyMap{
		"event": SuccessEvent,
		"path":  path,
	})
}

func (c *Communicator) SendSuccessFlushing(jobID int) error {
	return c.client.Trigger(fmt.Sprint(jobID), FlushedEvent, types.AnyMap{
		"event": SuccessEvent,
	})
}

func channelName(jobID int) string {
	return "cache-" + strconv.Itoa(jobID)
}
