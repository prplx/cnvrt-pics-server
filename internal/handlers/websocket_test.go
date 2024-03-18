package handlers_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/prplx/cnvrt/internal/handlers"
	"github.com/prplx/cnvrt/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandleWebsocket__should_add_remove_client_and_return_when_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mocks.NewMockWebsocketConnection(ctrl)
	communicator := mocks.NewMockCommunicator(ctrl)
	mocks := &Mocks{communicator: communicator}
	_, services := setup(t, mocks)

	err := errors.New("Some error")
	mockConn.EXPECT().Params("jobID").Return("1")
	communicator.EXPECT().AddClient(1, mockConn)
	mockConn.EXPECT().ReadMessage().Return(int64(1), []byte("test"), nil).Times(1)
	mockConn.EXPECT().ReadMessage().Return(int64(1), []byte("test"), err).Times(1)
	mockConn.EXPECT().Close()

	communicator.EXPECT().RemoveClient(1)

	h := handlers.NewHandlers(services)

	h.HandleWebsocket(mockConn)
}
