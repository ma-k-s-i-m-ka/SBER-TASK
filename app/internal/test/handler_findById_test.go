package test

import (
	"Sber/app/internal/cache"
	"Sber/app/internal/task"
	"Sber/app/internal/task/mocks"
	"Sber/app/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindByIdTask(t *testing.T) {
	router := httprouter.New()
	serviceMock := new(mocks.Service)
	handler := task.NewHandler(logger.GetLogger(), serviceMock, cache.NewCache())
	handler.Register(router)

	testCases := []struct {
		ID          int64
		ExpectedID  int64
		ExpectedErr error
	}{
		{ID: 1, ExpectedID: 1},
		{ID: 2, ExpectedID: 2},
		{ID: 3, ExpectedID: 3},
		{ID: 4, ExpectedID: 4},
		{ID: 5, ExpectedID: 5},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", fmt.Sprintf("/task/%d", testCase.ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		expectedTask := &task.Task{ID: testCase.ExpectedID}
		serviceMock.On("GetById", mock.Anything, testCase.ID).Return(expectedTask, nil).Once()
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		var responseTask task.Task
		err = json.NewDecoder(recorder.Body).Decode(&responseTask)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testCase.ExpectedID, responseTask.ID)
		serviceMock.AssertExpectations(t)
	}
}
