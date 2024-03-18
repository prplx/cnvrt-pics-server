package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prplx/cnvrt/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_HandleArchiveJob__should_call_achiver_service_when_all_conditions_are_met(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := int64(555)
	archiver := mocks.NewMockArchiver(ctrl)
	mocks := &Mocks{
		archiver: archiver,
	}
	archiver.EXPECT().Archive(jobID).Times(1)
	app, services := setup(t, mocks)

	r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/%d", archiveEndpoint, jobID), nil)

	resp, _ := app.Test(r, -1)

	assert.Equal(t, 200, resp.StatusCode)

	cleanUp(t, services)
}

func Test_HandleArchiveJob__should_return_500_when_jobID_is_not_a_number(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := "abc"
	mocks := &Mocks{}
	app, services := setup(t, mocks)

	r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", archiveEndpoint, jobID), nil)

	resp, _ := app.Test(r, -1)
	got, _ := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"errors":"jobID must be a number"}`, string(got))

	cleanUp(t, services)
}
