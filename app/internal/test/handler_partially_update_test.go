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

func TestPartiallyUpdateTask(t *testing.T) {
	router := httprouter.New()
	serviceMock := new(mocks.Service)
	handler := task.NewHandler(logger.GetLogger(), serviceMock, cache.NewCache())
	handler.Register(router)

	ptrString := func(s string) *string {
		return &s
	}

	ptrTime := func(t time.Time) *time.Time {
		return &t
	}

	ptrBool := func(b bool) *bool {
		return &b
	}

	testCases := []struct {
		ID          int64
		Input       task.PartiallyUpdateTask
		Expected    *task.Task
		ExpectedErr error
	}{
		{
			ID: 1,
			Input: task.PartiallyUpdateTask{
				ID:          1,
				Title:       ptrString("Обновленный заголовок"),
				Description: ptrString(""),
				Date:        ptrTime(time.Date(2023, 9, 20, 10, 0, 0, 0, time.UTC)),
				Status:      ptrBool(true),
			},
			Expected: &task.Task{
				ID:          1,
				Title:       "Обновленный заголовок",
				Description: " ",
				Date:        time.Date(2023, 9, 20, 10, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedErr: nil,
		},
		{
			ID: 2,
			Input: task.PartiallyUpdateTask{
				ID:          2,
				Title:       ptrString("Обновленный заголовок"),
				Description: ptrString("Новое описание"),
				Date:        ptrTime(time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC)),
				Status:      ptrBool(false),
			},
			Expected: &task.Task{
				ID:          2,
				Title:       "Обновленный заголовок",
				Description: "Новое описание",
				Date:        time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC),
				Status:      false,
			},
			ExpectedErr: nil,
		},
		{
			ID: 3,
			Input: task.PartiallyUpdateTask{
				ID:          3,
				Title:       ptrString(""),
				Description: ptrString("Новое описание"),
				Date:        ptrTime(time.Date(2023, 9, 23, 15, 0, 0, 0, time.UTC)),
				Status:      ptrBool(true),
			},
			Expected: &task.Task{
				ID:          3,
				Title:       "",
				Description: "Новое описание",
				Date:        time.Date(2023, 9, 23, 15, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedErr: nil,
		},
		{
			ID: 4,
			Input: task.PartiallyUpdateTask{
				ID:          4,
				Title:       ptrString("Обновленный заголовок"),
				Description: ptrString("Новое описание"),
				Date:        ptrTime(time.Date(2023, 9, 23, 15, 0, 0, 0, time.UTC)),
				Status:      ptrBool(true),
			},
			Expected: &task.Task{
				ID:          4,
				Title:       "Обновленный заголовок",
				Description: "Новое описание",
				Date:        time.Date(2023, 9, 23, 15, 0, 0, 0, time.UTC),
				Status:      true,
			},
			ExpectedErr: nil,
		},
		{
			ID: 5,
			Input: task.PartiallyUpdateTask{
				ID:          5,
				Title:       ptrString(""),
				Description: ptrString(""),
				Date:        ptrTime(time.Date(2023, 9, 23, 15, 0, 0, 0, time.UTC)),
				Status:      ptrBool(true),
			},
			Expected: &task.Task{
				ID:          5,
				Title:       "",
				Description: "",
				Date:        time.Date(2023, 9, 23, 15, 0, 0, 0, time.UTC),
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
		req, err := http.NewRequest("PATCH", fmt.Sprintf("/task/%d", testCase.ID), bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		serviceMock.On("PartiallyUpdate", mock.Anything, &testCase.Input).Return(testCase.Expected, testCase.ExpectedErr).Once()

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
