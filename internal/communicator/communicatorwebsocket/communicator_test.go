package communicatorwebsocket_test

import (
	"strconv"
	"testing"

	communicator "github.com/prplx/cnvrt/internal/communicator/communicatorwebsocket"
	"github.com/prplx/cnvrt/internal/mocks"
	"github.com/prplx/cnvrt/internal/types"
	"go.uber.org/mock/gomock"
)

func Test_NewCommunicator__should_send_message_from_the_cache_when_the_client_is_not_connected_yet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	fileID := int64(2)
	fileName := "test.jpg"
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation": communicator.ProcessingOperation,
		"event":     communicator.StartedEvent,
		"fileName":  fileName,
		"fileId":    fileID,
	})).Times(1).Return(nil)

	comm.SendStartProcessing(jobID, fileID, fileName)
	comm.AddClient(jobID, connectionMock)
}

func Test_NewCommunicator__should_send_message_from_the_cache_for_multiple_files_when_the_client_is_not_connected_yet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	operations := []struct {
		jobID    int64
		fileID   int64
		fileName string
	}{
		{jobID: jobID,
			fileID:   1,
			fileName: "test1.jpg"},
		{jobID: jobID,
			fileID:   2,
			fileName: "test2.jpg"},
		{jobID: jobID,
			fileID:   3,
			fileName: "test3.jpg"},
	}
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	for _, op := range operations {
		connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
			"operation": communicator.ProcessingOperation,
			"event":     communicator.StartedEvent,
			"fileName":  op.fileName,
			"fileId":    op.fileID,
		})).Times(1).Return(nil)
		comm.SendStartProcessing(op.jobID, op.fileID, op.fileName)
	}

	comm.AddClient(jobID, connectionMock)
}

func Test_NewCommunicator__should_send_start_processing_job_message_to_respective_connection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	fileID := int64(2)
	fileName := "test.jpg"
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation": communicator.ProcessingOperation,
		"event":     communicator.StartedEvent,
		"fileName":  fileName,
		"fileId":    fileID,
	})).Times(1).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.SendStartProcessing(jobID, fileID, fileName)
}

func Test_NewCommunicator__should_not_send_procession_job_message_if_connection_was_removed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	fileID := int64(2)
	fileName := "test.jpg"
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Any()).Times(0).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.RemoveClient(jobID)
	comm.SendStartProcessing(jobID, fileID, fileName)
}

func Test_NewCommunicator__should_send_error_processing_job_message_to_respective_connection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(2)
	fileID := int64(3)
	fileName := "test.jpg"
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation": communicator.ProcessingOperation,
		"event":     communicator.ErrorEvent,
		"fileName":  fileName,
		"fileId":    fileID,
	})).Times(1).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.SendErrorProcessing(jobID, fileID, fileName)
}

func Test_NewCommunicator__should_send_success_processing_job_to_respective_connection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(4)
	fileID := int64(5)
	sourceFileName := "test.jpg"
	targetFileName := "test.webp"

	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation":      communicator.ProcessingOperation,
		"event":          communicator.SuccessEvent,
		"fileId":         fileID,
		"sourceFile":     sourceFileName,
		"targetFile":     targetFileName,
		"sourceFileSize": strconv.FormatInt(100, 10),
		"targetFileSize": strconv.FormatInt(200, 10),
		"width":          100,
		"height":         200,
	})).Times(1).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.SendSuccessProcessing(jobID, types.SuccessResult{
		SourceFileName: sourceFileName,
		SourceFileID:   fileID,
		TargetFileName: targetFileName,
		SourceFileSize: 100,
		TargetFileSize: 200,
		Width:          100,
		Height:         200,
	})
}

func Test_NewCommunicator__should_send_start_archiving_message_for_respective_connection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation": communicator.ArchivingOperation,
		"event":     communicator.StartedEvent,
	})).Times(1).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.SendStartArchiving(jobID)
}

func Test_NewCommunicator__should_send_error_archiving_message_for_respective_connection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation": communicator.ArchivingOperation,
		"event":     communicator.ErrorEvent,
	})).Times(1).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.SendErrorArchiving(jobID)
}

func Test_NewCommunicator__should_send_success_archiving_message_for_respective_connection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	path := "/tmp/test.zip"
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation": communicator.ArchivingOperation,
		"event":     communicator.SuccessEvent,
		"path":      path,
	})).Times(1).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.SendSuccessArchiving(jobID, path)
}

func Test_NewCommunicator__should_send_success_flushing_message_for_respective_connection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(1)
	comm := communicator.NewCommunicator()
	connectionMock := mocks.NewMockWebsocketConnection(ctrl)

	connectionMock.EXPECT().WriteJSON(gomock.Eq(types.AnyMap{
		"operation": communicator.FlushingOperation,
		"event":     communicator.SuccessEvent,
	})).Times(1).Return(nil)

	comm.AddClient(jobID, connectionMock)
	comm.SendSuccessFlushing(jobID)
}
