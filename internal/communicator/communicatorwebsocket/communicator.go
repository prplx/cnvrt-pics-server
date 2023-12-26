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
	mu              sync.Mutex
	connections     map[int]types.WebsocketConnection
	processingCache map[int]*[]*types.AnyMap
}

func NewCommunicator() *Communicator {
	return &Communicator{
		mu:              sync.Mutex{},
		connections:     make(map[int]types.WebsocketConnection),
		processingCache: make(map[int]*[]*types.AnyMap),
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
	c.sendMessagesFromProcessingCache(jobID)
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
		return c.addMessageToProcessingCache(jobID, &message)
	}

	return conn.WriteJSON(message)
}

func (c *Communicator) SendErrorProcessing(jobID int, fileID int, fileName string) error {
	conn := c.connections[jobID]
	message := types.AnyMap{
		"operation": ProcessingOperation,
		"event":     ErrorEvent,
		"fileName":  fileName,
		"fileId":    fileID,
	}

	if conn == nil {
		return c.addMessageToProcessingCache(jobID, &message)
	}

	return conn.WriteJSON(message)
}

func (c *Communicator) SendSuccessProcessing(jobID int, result types.SuccessResult) error {
	conn := c.connections[jobID]
	message := types.AnyMap{
		"operation":      ProcessingOperation,
		"event":          SuccessEvent,
		"fileId":         result.SourceFileID,
		"sourceFile":     result.SourceFileName,
		"targetFile":     result.TargetFileName,
		"sourceFileSize": strconv.FormatInt(result.SourceFileSize, 10),
		"targetFileSize": strconv.FormatInt(result.TargetFileSize, 10),
		"width":          result.Width,
		"height":         result.Height,
	}

	if conn == nil {
		return c.addMessageToProcessingCache(jobID, &message)
	}

	return conn.WriteJSON(message)
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

func (c *Communicator) addMessageToProcessingCache(jobID int, message *types.AnyMap) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	slice := c.processingCache[jobID]
	if slice == nil {
		slice = &[]*types.AnyMap{}
	}
	*slice = append(*slice, message)
	c.processingCache[jobID] = slice

	return nil
}

func (c *Communicator) sendMessagesFromProcessingCache(jobID int) error {
	slice := c.processingCache[jobID]
	if slice == nil {
		return nil
	}

	conn := c.connections[jobID]
	if conn == nil {
		return nil
	}

	for _, message := range *slice {
		if err := conn.WriteJSON(*message); err != nil {
			return err
		}
	}

	delete(c.processingCache, jobID)
	return nil
}
