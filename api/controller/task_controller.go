package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

type TaskController interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetAllTasks(w http.ResponseWriter, r *http.Request)
	GetTasksofTeam(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
}

type taskController struct {
	taskService service.TaskService
}

func NewTaskController(taskService service.TaskService) TaskController {
	return taskController{
		taskService: taskService,
	}
}

// CreateTask creates a new task.
// @Summary Create New Task
// @Description CreateTask API is made for creating a new task in the task manager application.
// @Accept json
// @Produce json
// @Tags tasks
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param title formData string true "Title of the task (min length: 4, max length: 48)"
// @Param description formData string true "Description of the task (min length: 12, max length: 196)"
// @Param deadline formData time true "Deadline to Complete the task"
// @Param assigneeIndividual formData int64 false "ID of the individual assignee"
// @Param assigneeTeam formData int64 false "ID of the team assignee"
// @Param status formData string true "Status of the task (TO-DO, In-PROGRESS, COMPLETED, CLOSED)"
// @Param priority formData string true "Priority of the task (LOW, MEDIUM, HIGH, VERY HIGH)"
// @Success 200 {object} response.SuccessResponse "Task created successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request, either data is not valid or assignee privacy is Private."
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/tasks [post]
func (t taskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskToCreate request.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &taskToCreate)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadDataError, constant.EMPTY_STRING)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(taskToCreate)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	taskToCreate.CreatedAt = time.Now().UTC()
	taskToCreate.CreatedBy = userId
	taskId, err := t.taskService.CreateTask(taskToCreate)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	response := response.SuccessResponse{
		Code:    http.StatusText(http.StatusOK),
		Message: constant.TASK_CREATED,
		ID:      &taskId,
	}
	config.LoggerInstance.Info(constant.TASK_CREATED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// GetAllTasks fetches all tasks of user.
// @Summary Get all tasks
// @Description Get all tasks of user based on query parameters
// @Produce  json
// @Tags tasks
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param createdByMe query bool true "return tasks created by you if createdByMe set to true otherwise false."
// @Param limit query int false "Number of tasks to return per page (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Param search query string false "Search term to filter tasks"
// @Param status query string false "Filter tasks by status (TO-DO, In-PROGRESS, COMPLETED, CLOSED)"
// @Param sortByFilter query bool false "Sort tasks by create time (true for ascending, false for descending)"
// @Success 200 {object} []response.Task "Tasks fetched successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 422 {object} errorhandling.CustomError "Provide valid flag"
// @Failure 500 {object} errorhandling.CustomError "Internal server error"
// @Router /api/v1/tasks/{Flag} [get]
func (t taskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var taskQueryParams request.TaskQueryParams

	decoder := schema.NewDecoder()
	err := decoder.Decode(&taskQueryParams, r.URL.Query())
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.CreateCustomError(err.Error(), http.StatusText(http.StatusBadRequest)), constant.EMPTY_STRING)
		return
	}

	err = utils.Validate.Struct(taskQueryParams)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	if taskQueryParams.Limit == 0 {
		taskQueryParams.Limit = 10
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	tasks, err := t.taskService.GetAllTasks(userId, taskQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, tasks)
}

// @Summary Get all tasks of a team
// @Description Get all tasks of a team based on query parameters
// @Produce json
// @Tags tasks
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param TeamID path int64 true "Team ID"
// @Param limit query int false "Number of tasks to return per page (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Param search query string false "Search term to filter tasks"
// @Param status query string false "Filter tasks by status (TO-DO, In-PROGRESS, COMPLETED, CLOSED)"
// @Param sortByFilter query bool false "Sort tasks by create time (true for ascending, false for descending)"
// @Success 200 {object} []response.Task "Tasks fetched successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error"
// @Router /api/v1/tasks/team/{TeamID} [get]
// GetTasksOfTeam fetches all tasks of a specific team.
func (t taskController) GetTasksofTeam(w http.ResponseWriter, r *http.Request) {
	var taskQueryParams request.TaskQueryParams

	decoder := schema.NewDecoder()
	err := decoder.Decode(&taskQueryParams, r.URL.Query())
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadQueryParamsError, constant.EMPTY_STRING)
		return
	}

	err = utils.Validate.Struct(taskQueryParams)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	if taskQueryParams.Limit == 0 {
		taskQueryParams.Limit = 10
	}

	teamId, err := strconv.ParseInt(chi.URLParam(r, constant.TEAM_ID), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), constant.URL_PARAM_CONVERT_ERROR) {
			errorhandling.SendErrorResponse(r, w, errorhandling.ProvideValidParams, constant.EMPTY_STRING)
			return
		}
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	tasks, err := t.taskService.GetTasksofTeam(teamId, taskQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, tasks)
}

// UpdateTask updates a task.
// @Summary Update a task
// @Description Update a task based on provided parameters
// @Accept json
// @Produce json
// @Tags tasks
// @Param TaskID path int64 true "Task ID"
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param title formData string false "Title of the task (min length: 4, max length: 48)"
// @Param description formData string false "Description of the task (min length: 12, max length: 196)"
// @Param assigneeIndividual formData int64 false "ID of the individual assignee"
// @Param assigneeTeam formData int64 false "ID of the team assignee"
// @Param status formData string false "Status of the task (TO-DO, In-PROGRESS, COMPLETED, CLOSED)"
// @Param priority formData string false "Priority of the task (LOW, MEDIUM, HIGH, VERY HIGH)"
// @Success 200 {object} response.SuccessResponse "Task updated successfully"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 403 {object} errorhandling.CustomError "Not allowed to update task"
// @Failure 404 {object} errorhandling.CustomError "Task not found"
// @Failure 422 {object} errorhandling.CustomError "Task is closed"
// @Failure 500 {object} errorhandling.CustomError "Internal server error"
// @Router /api/v1/tasks/ [put]
func (t taskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var taskToUpdate request.UpdateTask

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &taskToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadDataError, constant.EMPTY_STRING)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	taskId, err := strconv.ParseInt(chi.URLParam(r, constant.TASK_ID), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), constant.URL_PARAM_CONVERT_ERROR) {
			errorhandling.SendErrorResponse(r, w, errorhandling.ProvideValidParams, constant.EMPTY_STRING)
			return
		}
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	taskToUpdate.ID = taskId

	err = utils.Validate.Struct(taskToUpdate)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	taskToUpdate.UpdatedBy = &userId
	taskToUpdate.UpdatedAt = new(time.Time)
	*taskToUpdate.UpdatedAt = time.Now().UTC()

	err = t.taskService.UpdateTask(taskToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	response := response.SuccessResponse{
		Code:    http.StatusText(http.StatusOK),
		Message: constant.TASK_UPDATED,
	}
	config.LoggerInstance.Info(constant.TASK_UPDATED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
