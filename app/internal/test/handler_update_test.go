package test

import (
	"Sber/app/internal/cache"
	"Sber/app/internal/task"
	"Sber/app/internal/task/mocks"
	"Sber/app/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUpdateTask(t *testing.T) {
	router := httprouter.New()
	serviceMock := new(mocks.Service)
	handler := task.NewHandler(logger.GetLogger(), serviceMock, cache.NewCache())
	handler.Register(router)

	testCases := []struct {
		ID          int64
		Input       task.Task
		Expected    *task.Task
		ExpectedErr error
	}{
		{
			ID: 1,
			Input: task.Task{
				ID:          1,
				Title:       "Обновленная задача 1",
				Description: "Новое описание 1",
				Date:        time.Date(2023, 9, 20, 10, 0, 0, 0, time.UTC),
				Status:      true,
			},
			Expected: &task.Task{
				ID:          1,
				Title:       "Обновленная задача 1",
				Description: "Новое описание 1",
				Date:        time.Date(2023, 9, 20, 10, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedErr: nil,
		},
		{
			ID: 2,
			Input: task.Task{
				ID:          2,
				Title:       "Обновленная задача 2",
				Description: "Новое описание 2",
				Date:        time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				Status:      false,
			},
			Expected: &task.Task{
				ID:          2,
				Title:       "Обновленная задача 2",
				Description: "Новое описание 2",
				Date:        time.Date(2023, 9, 21, 12, 0, 0, 0, time.UTC),
				Status:      false,
			},
			ExpectedErr: nil,
		},
		{
			ID: 3,
			Input: task.Task{
				ID:          3,
				Title:       "Обновленная задача 3",
				Description: "Новое описание 3",
				Date:        time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC),
				Status:      true,
			},
			Expected: &task.Task{
				ID:          3,
				Title:       "Обновленная задача 3",
				Description: "Новое описание 3",
				Date:        time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedErr: nil,
		},
		{
			ID: 4,
			Input: task.Task{
				ID:          4,
				Title:       "Обновленная задача 4",
				Description: "Новое описание 4",
				Date:        time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC),
				Status:      true,
			},
			Expected: &task.Task{
				ID:          4,
				Title:       "Обновленная задача 4",
				Description: "Новое описание 4",
				Date:        time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedErr: nil,
		},
		{
			ID: 5,
			Input: task.Task{
				ID:          5,
				Title:       "Обновленная задача 5",
				Description: "Новое описание 5",
				Date:        time.Date(2023, 9, 24, 14, 0, 0, 0, time.UTC),
				Status:      true,
			},
			Expected: &task.Task{
				ID:          5,
				Title:       "Обновленная задача 5",
				Description: "Новое описание 5",
				Date:        time.Date(2023, 9, 24, 14, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		requestBody, err := json.Marshal(testCase.Input)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("PUT", fmt.Sprintf("/task/%d", testCase.ID), bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		serviceMock.On("Update", mock.Anything, &testCase.Input).Return(testCase.Expected, testCase.ExpectedErr).Once()

		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)

		if testCase.ExpectedErr == nil {
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			var updatedTask task.Task
			err = json.NewDecoder(recorder.Body).Decode(&updatedTask)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.Expected, &updatedTask)
		} else {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			assert.True(t, strings.Contains(recorder.Body.String(), "error"))
		}
		serviceMock.AssertExpectations(t)
	}
}
