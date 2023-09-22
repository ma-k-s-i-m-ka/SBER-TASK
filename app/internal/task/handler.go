package task

import (
	"Sber/app/internal/apperror"
	"Sber/app/internal/cache"
	"Sber/app/internal/handler"
	"Sber/app/internal/model"
	"Sber/app/internal/response"
	"Sber/app/pkg/logger"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sort"
)

const (
	taskURL          = "/task"
	taskAllURL       = "/task_all"
	taskAllStatusURL = "/task_all_status"
	taskAvailableURL = "/task_all_available"
	taskIdURL        = "/task/:id"
)

type Handler struct {
	log         logger.Logger
	taskService Service
	cache       *cache.Cache
}

func NewHandler(log logger.Logger, taskService Service, cache *cache.Cache) handler.Hand {
	return &Handler{
		log:         log,
		taskService: taskService,
		cache:       cache,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, taskIdURL, h.GetTaskById)
	router.HandlerFunc(http.MethodGet, taskAllURL, h.FindAllTasks)
	router.HandlerFunc(http.MethodPost, taskAllStatusURL, h.FindAllStatusTasks)
	router.HandlerFunc(http.MethodPost, taskAvailableURL, h.FindDateAllAvailableTask)
	router.HandlerFunc(http.MethodPost, taskURL, h.CreateTask)
	router.HandlerFunc(http.MethodPut, taskIdURL, h.UpdateTask)
	router.HandlerFunc(http.MethodPatch, taskIdURL, h.PartiallyUpdateTask)
	router.HandlerFunc(http.MethodDelete, taskIdURL, h.DeleteTask)
}

// @Summary Создать задачу
// @Description Создает новую задачу
// @Accept  json
// @Produce  json
// @Example CreateTask
//
//	{
//	 "title": "Новая задача",
//
// / "description": "Описание новой задачи",
//
//	 "date": "2023-09-22T09:00:00Z",
//	 "status": false
//	}
//
// @Param input body CreateTask true "Данные для создания задачи"
// @Success 201 {object} Task
// @Router /task [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: CREATE TASK")

	var input CreateTask

	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	h.log.Info("Input: ", input)

	task, err := h.taskService.Create(r.Context(), &input)
	if err != nil {
		response.InternalError(w, fmt.Sprintf("cannot create task: %v", err), "")
		return
	}

	response.JSON(w, http.StatusCreated, task)
}

// @Summary Получить задачу по идентификатору
// @Description Получает задачу по заданному идентификатору
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Success 200 {object} Task
// @Router /task/{id} [get]
func (h *Handler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: GET TASK BY ID")

	id, err := handler.ReadIdParam64(r)
	h.log.Info("Input: ", id)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	cacheTask, ok := h.cache.Task[id]
	if ok {
		h.log.Info("GOT TASK FROM CACHE BY ID")
		response.JSON(w, http.StatusOK, cacheTask)
		return
	}

	task, err := h.taskService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		h.log.Error(err)
		response.InternalError(w, err.Error(), "")
		return
	}
	response.JSON(w, http.StatusOK, task)
}

// @Summary Получить все задачи
// @Description Получает список всех задач
// @Accept json
// @Produce json
// @Success 200 {array} Task
// @Router /tasks [get]
func (h *Handler) FindAllTasks(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: GET ALL TASKS")

	cacheTasks := make([]*model.Task, 0)
	for _, task := range h.cache.Task {
		cacheTasks = append(cacheTasks, task)
	}
	if len(cacheTasks) > 0 {
		sort.Slice(cacheTasks, func(i, j int) bool {
			return cacheTasks[i].ID < cacheTasks[j].ID
		})
		h.log.Info("GOT TASKS FROM CACHE")
		response.JSON(w, http.StatusOK, cacheTasks)
		return
	}
	tasks, err := h.taskService.FindAll()
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	if len(*tasks) == 0 {
		response.NotFound(w)
		return
	}

	response.JSON(w, http.StatusOK, tasks)
}

// @Summary Получить все задачи с определенным статусом
// @Description Получает список всех задач с заданным статусом
// @Accept json
// @Produce json
// @Param status query string true "Запрос на получение задач с определенным статусом"
// @Success 200 {array} Task
// @Router /tasks/status [post]
func (h *Handler) FindAllStatusTasks(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: GET ALL AVAILABLE STATUS TASKS")
	var input Task
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	status := input.Status
	cachedTasks := make([]*model.Task, 0)
	for _, task := range h.cache.Task {
		if task.Status == status {
			cachedTasks = append(cachedTasks, task)
		}
	}
	if len(cachedTasks) > 0 {
		sort.Slice(cachedTasks, func(i, j int) bool {
			return cachedTasks[i].ID < cachedTasks[j].ID
		})
		h.log.Info("GOT STATUS TASKS FROM CACHE")
		response.JSON(w, http.StatusOK, cachedTasks)
		return
	}

	tasks, err := h.taskService.FindAllStatus(r.Context(), status)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}
	if len(*tasks) == 0 {
		response.NotFound(w)
		return
	}
	response.JSON(w, http.StatusOK, tasks)
}

// @Summary Получить все задачи по определенной дате
// @Description Получает список всех доступных задач по заданной дате
// @Accept json
// @Produce json
// @Param date query string true "Запрос на получение задач с определенной датой и статусом"
// @Success 200 {array} Task
// @Router /tasks/date [post]
func (h *Handler) FindDateAllAvailableTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: GET ALL AVAILABLE DATE TASKS")
	var input Task
	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}
	date := input.Date
	status := input.Status

	cachedTasks := make([]*model.Task, 0)
	for _, task := range h.cache.Task {
		if task.Date == date && task.Status == status {
			cachedTasks = append(cachedTasks, task)
		}
	}
	if len(cachedTasks) > 0 {
		sort.Slice(cachedTasks, func(i, j int) bool {
			return cachedTasks[i].ID < cachedTasks[j].ID
		})
		h.log.Info("GOT STATUS TASKS FROM CACHE")
		response.JSON(w, http.StatusOK, cachedTasks)
		return
	}

	tasks, err := h.taskService.FindDateAllAvailable(r.Context(), date, status)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	if len(*tasks) == 0 {
		response.NotFound(w)
		return
	}
	response.JSON(w, http.StatusOK, tasks)
}

// @Summary Обновить задачу
// @Description Обновляет существующую задачу
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Param input body Task true "Данные для обновления задачи"
// @Success 200 {object} Task
// @Router /task/{id} [put]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: UPDATE TASK")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input Task
	input.ID = id

	if err = response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	task, err := h.taskService.Update(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	response.JSON(w, http.StatusOK, task)
}

// @Summary Частичное обновление задачи
// @Description Обновляет часть данных существующей задачи
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Param input body PartiallyUpdateTask true "Часть данных для обновления задачи, строки могут быть пустыми"
// @Success 200 {object} Task
// @Router /task/{id} [patch]
func (h *Handler) PartiallyUpdateTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: PARTIALLY UPDATE TASK")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	var input PartiallyUpdateTask
	input.ID = id

	if err = response.ReadJSON(w, r, &input); err != nil {
		response.BadRequest(w, err.Error(), apperror.ErrInvalidRequestBody.Error())
		return
	}

	task, err := h.taskService.PartiallyUpdate(r.Context(), &input)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
	}
	response.JSON(w, http.StatusOK, task)
}

// @Summary Удалить задачу
// @Description Удаляет задачу по заданному идентификатору
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Success 200 {string} string
// @Router /task/{id} [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HANDLER: DELETE TASK")

	id, err := handler.ReadIdParam64(r)
	if err != nil {
		response.BadRequest(w, err.Error(), "")
		return
	}

	err = h.taskService.Delete(id)
	if err != nil {
		if errors.Is(err, apperror.ErrEmptyString) {
			response.NotFound(w)
			return
		}
		response.InternalError(w, err.Error(), "wrong on the server")
		return
	}
	response.JSON(w, http.StatusOK, "TASK DELETED")
}
