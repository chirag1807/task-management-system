package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/go-chi/chi/v5"
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

func (t taskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.TitleKey:              "string|minLen:4|maxLen:24|required",
		constant.DescriptionKey:        "string|minLen:12|maxLen:108|required",
		constant.DeadlineKey:           "required",
		constant.AssigneeIndividualKey: `number`,
		constant.AssigneeTeamKey:       "number",
		constant.StatusKey:             "string|in:TO-DO,In-Progress,Completed,Closed|required",
		constant.PriorityKey:           "string|in:Low,Medium,High,Very High|required",
	}
	var taskToCreate request.Task

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &taskToCreate, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &taskToCreate)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	taskToCreate.CreatedAt = time.Now().UTC()
	taskToCreate.CreatedBy = userId
	taskId, err := t.taskService.CreateTask(taskToCreate)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	response := response.SuccessResponse{
		Message: constant.TASK_CREATED,
		ID:      &taskId,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t taskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var queryParams = map[string]string{
		constant.LimitKey:        "number|default:10",
		constant.OffsetKey:       "number|default:0",
		constant.SearchKey:       "string",
		constant.StatusFilterKey: "string|in:TO-DO,In-Progress,Completed,Closed",
		constant.SortByFilterKey: "bool",
	}
	var queryParamFilters = map[string]string{
		constant.LimitKey:        "int",
		constant.OffsetKey:       "int",
		constant.SortByFilterKey: "bool",
	}

	var taskQueryParams request.TaskQueryParams

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &taskQueryParams, nil, nil, &queryParams, &queryParamFilters, nil)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	flag, err := strconv.Atoi(chi.URLParam(r, "Flag"))
	if err != nil {
		if strings.Contains(err.Error(), "strconv.Atoi: parsing") {
			errorhandling.SendErrorResponse(w, errorhandling.ProvideValidFlag)
			return
		}
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if flag == 0 || flag == 1 {
		tasks, err := t.taskService.GetAllTasks(userId, flag, taskQueryParams)
		if err != nil {
			errorhandling.SendErrorResponse(w, err)
			return
		}
		response := response.Tasks{
			Tasks: tasks,
		}
		utils.SendSuccessResponse(w, http.StatusOK, response)
	} else {
		errorhandling.SendErrorResponse(w, errorhandling.ProvideValidFlag)
	}
}

func (t taskController) GetTasksofTeam(w http.ResponseWriter, r *http.Request) {
	var queryParams = map[string]string{
		constant.LimitKey:        "number|default:10",
		constant.OffsetKey:       "number|default:0",
		constant.SearchKey:       "string",
		constant.StatusFilterKey: "string|in:TO-DO,In-Progress,Completed,Closed",
		constant.SortByFilterKey: "bool",
	}
	var queryParamFilters = map[string]string{
		constant.LimitKey:        "int",
		constant.OffsetKey:       "int",
		constant.SortByFilterKey: "bool",
	}

	var taskQueryParams request.TaskQueryParams

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &taskQueryParams, nil, nil, &queryParams, &queryParamFilters, nil)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

	teamID, err := strconv.ParseInt(chi.URLParam(r, "TeamID"), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), "strconv.Atoi: parsing") {
			errorhandling.SendErrorResponse(w, errorhandling.ProvideValidParams)
			return
		}
		errorhandling.SendErrorResponse(w, err)
		return
	}

	tasks, err := t.taskService.GetTasksofTeam(teamID, taskQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.Tasks{
		Tasks: tasks,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t taskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.TaskIdKey:             "number|required",
		constant.TitleKey:              "string|minLen:4|maxLen:24",
		constant.DescriptionKey:        "string|minLen:12|maxLen:108",
		constant.AssigneeIndividualKey: `number`,
		constant.AssigneeTeamKey:       "number",
		constant.StatusKey:             "string|in:TO-DO,In-Progress,Completed,Closed",
		constant.PriorityKey:           "string|in:Low,Medium,High,Very High",
	}
	var taskToUpdate request.Task

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &taskToUpdate, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &taskToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	taskToUpdate.UpdatedBy = &userId
	taskToUpdate.UpdatedAt = new(time.Time)
	*taskToUpdate.UpdatedAt = time.Now().UTC()

	err = t.taskService.UpdateTask(taskToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	response := response.SuccessResponse{
		Message: constant.TASK_UPDATED,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)

}
