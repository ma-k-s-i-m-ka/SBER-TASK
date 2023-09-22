package test

import (
	"Sber/app/internal/cache"
	"Sber/app/internal/task"
	"Sber/app/internal/task/mocks"
	"Sber/app/pkg/logger"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteTask(t *testing.T) {
	router := httprouter.New()
	serviceMock := new(mocks.Service)
	handler := task.NewHandler(logger.GetLogger(), serviceMock, cache.NewCache())
	handler.Register(router)

	testCases := []struct {
		ID          int64
		ExpectedErr error
	}{
		{ID: 1, ExpectedErr: nil},
		{ID: 2, ExpectedErr: nil},
		{ID: 3, ExpectedErr: nil},
		{ID: 4, ExpectedErr: nil},
		{ID: 5, ExpectedErr: nil},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/task/%d", testCase.ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		serviceMock.On("Delete", mock.AnythingOfType("int64")).Return(testCase.ExpectedErr).Once()
		router.ServeHTTP(recorder, req)
		if testCase.ExpectedErr == nil {
			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, "TASK DELETED", strings.Trim(recorder.Body.String(), "\""))
		} else {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		}
		serviceMock.AssertExpectations(t)
	}
}
