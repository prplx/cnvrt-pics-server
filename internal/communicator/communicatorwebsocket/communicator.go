package communicatorwebsocket

import (
	"strconv"
	"sync"

	"github.com/prplx/cnvrt/internal/types"
)

const (
	ProcessingOperation = "processing"
	ArchivingOperation  = "archiving"
	FlushingOperation   = "flushing"
	StartedEvent        = "started"
	ErrorEvent          = "error"
	SuccessEvent        = "success"
)

type Communicator struct {
	mu                   sync.Mutex
	connections          map[int]types.WebsocketConnection
	startProcessingCache map[int]types.AnyMap
}

func NewCommunicator() *Communicator {
	return &Communicator{
		mu:                   sync.Mutex{},
		connections:          make(map[int]types.WebsocketConnection),
		startProcessingCache: make(map[int]types.AnyMap),
	}
}

func (c *Communicator) AddClient(jobID int, connection types.WebsocketConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn := c.connections[jobID]
	if conn != nil {
		return
	}
	c.connections[jobID] = connection
	message := c.startProcessingCache[jobID]
	if message != nil {
		connection.WriteJSON(message)
		delete(c.startProcessingCache, jobID)
	}
}

func (c *Communicator) RemoveClient(jobID int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.connections, jobID)
}

func (c *Communicator) SendStartProcessing(jobID int, fileID int, fileName string) error {
	conn := c.connections[jobID]
	message := types.AnyMap{
		"operation": ProcessingOperation,
		"event":     StartedEvent,
		"fileName":  fileName,
		"fileId":    fileID,
	}

	if conn == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.startProcessingCache[jobID] = message
		return nil
	}

	return conn.WriteJSON(message)
}

func (c *Communicator) SendErrorProcessing(jobID int, fileID int, fileName string) error {
	conn := c.connections[jobID]
	if conn == nil {
		return nil
	}

	return conn.WriteJSON(types.AnyMap{
		"operation": ProcessingOperation,
		"event":     ErrorEvent,
		"fileName":  fileName,
		"fileId":    fileID,
	})
}

func (c *Communicator) SendSuccessProcessing(jobID int, result types.SuccessResult) error {
	conn := c.connections[jobID]
	if conn == nil {
		return nil
	}

	return conn.WriteJSON(types.AnyMap{
		"operation":      ProcessingOperation,
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
	conn := c.connections[jobID]
	if conn == nil {
		return nil
	}

	return conn.WriteJSON(types.AnyMap{
		"operation": ArchivingOperation,
		"event":     StartedEvent,
	})
}

func (c *Communicator) SendErrorArchiving(jobID int) error {
	conn := c.connections[jobID]
	if conn == nil {
		return nil
	}

	return conn.WriteJSON(types.AnyMap{
		"operation": ArchivingOperation,
		"event":     ErrorEvent,
	})
}

func (c *Communicator) SendSuccessArchiving(jobID int, path string) error {
	conn := c.connections[jobID]
	if conn == nil {
		return nil
	}

	return conn.WriteJSON(types.AnyMap{
		"operation": ArchivingOperation,
		"event":     SuccessEvent,
		"path":      path,
	})
}

func (c *Communicator) SendSuccessFlushing(jobID int) error {
	conn := c.connections[jobID]
	if conn == nil {
		return nil
	}

	return conn.WriteJSON(types.AnyMap{
		"operation": FlushingOperation,
		"event":     SuccessEvent,
	})
}
