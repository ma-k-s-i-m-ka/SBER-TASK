package task

import (
	"Sber/app/internal/apperror"
	"Sber/app/internal/cache"
	"Sber/app/internal/model"
	"Sber/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
	"time"
)

var _ Storage = &TaskStorage{}

/// Структура DoctorStorage содержащая поля для работы с БД \\\

type TaskStorage struct {
	log            logger.Logger
	conn           *pgx.Conn
	requestTimeout time.Duration
	cache          *cache.Cache
}

func NewStorage(storage *pgx.Conn, requestTimeout int, cache *cache.Cache) Storage {
	return &TaskStorage{
		log:            logger.GetLogger(),
		conn:           storage,
		requestTimeout: time.Duration(requestTimeout) * time.Second,
		cache:          cache,
	}
}

func (d *TaskStorage) Create(task *Task) (*Task, error) {
	d.log.Info("POSTGRES: CREATE TASK")

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`INSERT INTO Task (title, description, date, status)
			 VALUES($1,$2,$3,$4) 
			 RETURNING id`,
		task.Title, task.Description, task.Date, task.Status)

	err := row.Scan(&task.ID)
	if err != nil {
		err = fmt.Errorf("failed to execute create task query: %v", err)
		d.log.Error(err)
		return nil, err
	}

	modelDelivery := &model.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Date:        task.Date,
		Status:      task.Status,
	}

	d.cache.Task[task.ID] = modelDelivery

	return task, nil
}

func (d *TaskStorage) FindById(id int64) (*Task, error) {
	d.log.Info("POSTGRES: GET TASK BY ID")

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row := d.conn.QueryRow(ctx,
		`SELECT * FROM Task
			 WHERE id = $1`, id)

	task := &Task{}
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute find task by id query: %v", err)
		d.log.Error(err)
		return nil, err
	}
	return task, nil
}

func (d *TaskStorage) FindAll() ([]Task, error) {
	d.log.Info("POSTGRES: GET ALL TASKS")

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	rows, err := d.conn.Query(ctx,
		`SELECT * FROM Task`)
	if err != nil {
		err = fmt.Errorf("failed to SELLECT: %v", err)
		d.log.Error(err)
		return nil, err
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apperror.ErrEmptyString
			}
			err = fmt.Errorf("failed to execute find all tasks query: %v", err)
			d.log.Error(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (d *TaskStorage) FindAllStatus(status bool) ([]Task, error) {
	d.log.Info("POSTGRES: GET ALL AVAILABLE STATUS TASKS")

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	rows, err := d.conn.Query(ctx,
		`SELECT * FROM Task
			 WHERE status=$1`,
		status)
	if err != nil {
		err = fmt.Errorf("failed to SELLECT: %v", err)
		d.log.Error(err)
		return nil, err
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apperror.ErrEmptyString
			}
			err = fmt.Errorf("failed to execute find all available status tasks query: %v", err)
			d.log.Error(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (d *TaskStorage) FindDateAllAvailable(date time.Time, status bool) ([]Task, error) {
	d.log.Info("POSTGRES: GET ALL AVAILABLE DATE TASKS")

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	rows, err := d.conn.Query(ctx,
		`SELECT * FROM Task
			 WHERE date=$1 AND status=$2`,
		date, status)
	if err != nil {
		err = fmt.Errorf("failed to SELLECT: %v", err)
		d.log.Error(err)
		return nil, err
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apperror.ErrEmptyString
			}
			err = fmt.Errorf("failed to execute find all available tasks query: %v", err)
			d.log.Error(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (d *TaskStorage) Update(task *Task) (*Task, error) {
	d.log.Info("POSTGRES: UPDATE TASK")

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	row, err := d.conn.Exec(ctx,
		`UPDATE Task
			SET title=$1, description=$2, date=$3, status=$4
			WHERE id =$5`,
		task.Title, task.Description, task.Date, task.Status, task.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrEmptyString
		}
		err = fmt.Errorf("failed to execute update task query: %v", err)
		d.log.Error(err)
		return nil, err
	}

	if row.RowsAffected() == 0 {
		return nil, apperror.ErrEmptyString
	}

	if _, exists := d.cache.Task[task.ID]; exists {
		d.cache.Task[task.ID] = &model.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Date:        task.Date,
			Status:      task.Status,
		}
	}
	return task, nil
}

func (d *TaskStorage) PartiallyUpdate(task *PartiallyUpdateTask) (*Task, error) {
	d.log.Info("POSTGRES: PARTIALLY UPDATE TASK")

	values := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if task.Title != nil {
		values = append(values, fmt.Sprintf("title=$%d", argId))
		args = append(args, *task.Title)
		argId++
	}
	if task.Description != nil {
		values = append(values, fmt.Sprintf("description=$%d", argId))
		args = append(args, *task.Description)
		argId++
	}
	if task.Date != nil {
		values = append(values, fmt.Sprintf("date=$%d", argId))
		args = append(args, *task.Date)
		argId++
	}
	if task.Status != nil {
		values = append(values, fmt.Sprintf("status=$%d", argId))
		args = append(args, *task.Status)
		argId++
	}

	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE Task  SET %s WHERE id = $%d", valuesQuery, argId)
	args = append(args, task.ID)

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update task partially: %v", err)
	}

	if result.RowsAffected() == 0 {
		return nil, apperror.ErrEmptyString
	}

	updatedTask := &Task{
		ID: task.ID,
	}

	if _, exists := d.cache.Task[task.ID]; exists {
		cachedTask := d.cache.Task[task.ID]
		if cachedTask != nil {
			if task.Title != nil {
				updatedTask.Title = *task.Title
				cachedTask.Title = *task.Title
			} else {
				updatedTask.Title = cachedTask.Title
			}
			if task.Description != nil {
				updatedTask.Description = *task.Description
				cachedTask.Description = *task.Description
			} else {
				updatedTask.Description = cachedTask.Description
			}
			if task.Date != nil {
				updatedTask.Date = *task.Date
				cachedTask.Date = *task.Date
			} else {
				updatedTask.Date = cachedTask.Date
			}
			if task.Status != nil {
				updatedTask.Status = *task.Status
				cachedTask.Status = *task.Status
			} else {
				updatedTask.Status = cachedTask.Status
			}
		}
	}

	return updatedTask, nil
}

func (d *TaskStorage) Delete(id int64) error {
	d.log.Info("POSTGRES: DELETE TASK")

	ctx, cancel := context.WithTimeout(context.Background(), d.requestTimeout)
	defer cancel()

	result, err := d.conn.Exec(ctx,
		`DELETE FROM Task WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrEmptyString
	}

	if _, exists := d.cache.Task[id]; exists {
		delete(d.cache.Task, id)
	}

	return nil
}

func CacheForTask(dbConn *pgx.Conn, cache *cache.Cache) error {
	rows, err := dbConn.Query(context.Background(), `SELECT * FROM Task`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Date, &task.Status)
		if err != nil {
			return err
		}
		modelTask := &model.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Date:        task.Date,
			Status:      task.Status,
		}
		cache.Task[task.ID] = modelTask
	}
	return nil
}
