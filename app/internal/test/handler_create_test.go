package test

import (
	"Sber/app/internal/cache"
	"Sber/app/internal/task"
	"Sber/app/internal/task/mocks"
	"Sber/app/pkg/logger"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	router := httprouter.New()
	serviceMock := new(mocks.Service)
	handler := task.NewHandler(logger.GetLogger(), serviceMock, cache.NewCache())
	handler.Register(router)

	testCases := []struct {
		Input       task.CreateTask
		ExpectedID  int64
		ExpectedErr error
	}{
		{
			Input: task.CreateTask{
				Title:       "Заголовок задачи 1",
				Description: "Описание задачи 1",
				Date:        time.Date(2023, 9, 20, 10, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedID:  1,
			ExpectedErr: nil,
		},
		{
			Input: task.CreateTask{
				Title:       "Заголовок задачи 2",
				Description: "Описание задачи 2",
				Date:        time.Date(2023, 9, 21, 11, 0, 0, 0, time.UTC),
				Status:      false,
			},
			ExpectedID:  2,
			ExpectedErr: nil,
		},
		{
			Input: task.CreateTask{
				Title:       "Заголовок задачи 3",
				Description: "Описание задачи 3",
				Date:        time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedID:  3,
			ExpectedErr: nil,
		},
		{
			Input: task.CreateTask{
				Title:       "Заголовок задачи 4",
				Description: "Описание задачи 4",
				Date:        time.Date(2023, 9, 23, 13, 0, 0, 0, time.UTC),
				Status:      false,
			},
			ExpectedID:  4,
			ExpectedErr: nil,
		},
		{
			Input: task.CreateTask{
				Title:       "Заголовок задачи 5",
				Description: "Описание задачи 5",
				Date:        time.Date(2023, 9, 24, 14, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedID:  5,
			ExpectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		requestBody, err := json.Marshal(testCase.Input)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/task", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()

		expectedTask := &task.Task{
			ID:          testCase.ExpectedID,
			Title:       testCase.Input.Title,
			Description: testCase.Input.Description,
			Date:        testCase.Input.Date,
			Status:      testCase.Input.Status,
		}
		serviceMock.On("Create", mock.Anything, &testCase.Input).Return(expectedTask, testCase.ExpectedErr).Once()
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		var createdTask task.Task
		err = json.NewDecoder(recorder.Body).Decode(&createdTask)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedTask, &createdTask)
		serviceMock.AssertExpectations(t)
	}
}
