package task

import (
	"Sber/app/internal/apperror"
	"Sber/app/pkg/logger"
	"context"
	"errors"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Service
type Service interface {
	Create(ctx context.Context, task *CreateTask) (*Task, error)
	GetById(ctx context.Context, id int64) (*Task, error)
	FindAll() (*[]Task, error)
	FindAllStatus(ctx context.Context, status bool) (*[]Task, error)
	FindDateAllAvailable(ctx context.Context, date time.Time, status bool) (*[]Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	PartiallyUpdate(ctx context.Context, task *PartiallyUpdateTask) (*Task, error)
	Delete(id int64) error
}

type service struct {
	log     logger.Logger
	storage Storage
}

func NewService(storage Storage, log logger.Logger) Service {
	return &service{
		log:     log,
		storage: storage,
	}
}

func (s *service) Create(ctx context.Context, input *CreateTask) (*Task, error) {
	s.log.Info("SERVICE: CREATE TASK")

	t := Task{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Status:      input.Status,
	}

	task, err := s.storage.Create(&t)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *service) GetById(ctx context.Context, id int64) (*Task, error) {
	s.log.Info("SERVICE: GET TASK BY ID")

	task, err := s.storage.FindById(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.log.Warn("cannot find task by id:", err)
		return nil, err
	}
	return task, nil
}

func (s *service) FindAll() (*[]Task, error) {
	s.log.Info("SERVICE: GET ALL DOCTORS")

	task, err := s.storage.FindAll()
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.log.Warnf("cannot find tasks by id: %v", err)
		return nil, err
	}
	return &task, nil
}

func (s *service) FindAllStatus(ctx context.Context, status bool) (*[]Task, error) {
	s.log.Info("SERVICE: GET ALL AVAILABLE STATUS TASKS")

	tasks, err := s.storage.FindAllStatus(status)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.log.Warnf("cannot find available status tasks by id: %v", err)
		return nil, err
	}
	return &tasks, nil
}

func (s *service) FindDateAllAvailable(ctx context.Context, date time.Time, status bool) (*[]Task, error) {
	s.log.Info("SERVICE: GET ALL AVAILABLE DATE TASKS")

	tasks, err := s.storage.FindDateAllAvailable(date, status)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			return nil, err
		}
		s.log.Warnf("cannot find available tasks by id: %v", err)
		return nil, err
	}
	return &tasks, nil
}

func (s *service) Update(ctx context.Context, task *Task) (*Task, error) {
	s.log.Info("SERVICE: UPDATE TASK")

	_, err := s.storage.FindById(task.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.log.Errorf("failed to get task: %v", err)
		}
		return nil, err
	}

	task, err = s.storage.Update(task)
	if err != nil {
		s.log.Errorf("failed to update task: %v", err)
		return nil, err
	}
	return task, nil
}

func (s *service) PartiallyUpdate(ctx context.Context, task *PartiallyUpdateTask) (*Task, error) {
	s.log.Info("SERVICE: PARTIALLY UPDATE TASK")

	_, err := s.storage.FindById(task.ID)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.log.Errorf("failed to get task: %v", err)
		}
		return nil, err
	}

	updatedTask, err := s.storage.PartiallyUpdate(task)
	if err != nil {
		s.log.Errorf("failed to partially update task: %v", err)
		return nil, err
	}
	return updatedTask, nil
}

func (s *service) Delete(id int64) error {
	s.log.Info("SERVICE: DELETE TASK")

	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, apperror.ErrEmptyString) {
			s.log.Warn("failed to delete task:", err)
		}
		return err
	}
	return nil
}
