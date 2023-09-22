package task

import "time"

type Storage interface {
	Create(task *Task) (*Task, error)
	FindById(id int64) (*Task, error)
	FindAll() ([]Task, error)
	FindAllStatus(status bool) ([]Task, error)
	FindDateAllAvailable(date time.Time, status bool) ([]Task, error)
	Update(task *Task) (*Task, error)
	PartiallyUpdate(task *PartiallyUpdateTask) (*Task, error)
	Delete(id int64) error
}
