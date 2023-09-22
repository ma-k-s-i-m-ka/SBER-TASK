package task

import "time"

// @Example Task
// {
// "id": 1, "title": "Задача 1",
// "description": "Описание задачи 1",
// "date": "2023-09-21T12:00:00Z",
// "status": false
// }
type Task struct {
	ID          int64     `json:"id" example:"1"`
	Title       string    `json:"title" example:"Задача 1"`
	Description string    `json:"description"example:"Описание задачи 1"`
	Date        time.Time `json:"date"example:"2023-09-21T12:00:00Z"`
	Status      bool      `json:"status"example:"false"`
}

// @Example CreateTask
// {
// "title": "Новая задача",
// "description": "Описание новой задачи",
// "date": "2023-09-22T09:00:00Z",
// "status": false
// }
type CreateTask struct {
	Title       string    `json:"title" example:"Новая задача"`
	Description string    `json:"description" example:"Описание новой задачи"`
	Date        time.Time `json:"date" example:"2023-09-22T09:00:00Z"`
	Status      bool      `json:"status" example:"false"`
}

// @Example PartiallyUpdateTask
// {
// "id": 1,
// "title": "Обновленный заголовок (Может быть пустым)",
// "description": "Обновленное описание (Может быть пустым)",
// "date": "2023-09-23T14:00:00Z (Может быть пустым)",
// "status": true (Может быть пустым)
// }
type PartiallyUpdateTask struct {
	ID          int64      `json:"id" example:"1"`
	Title       *string    `json:"title" example:"Обновленная Задача 1"`
	Description *string    `json:"description" example:"Обновленное Описание задачи 1"`
	Date        *time.Time `json:"date" example:"Обновленная дата 2023-09-21T12:00:00Z"`
	Status      *bool      `json:"status" example:"false"`
}
