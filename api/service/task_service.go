package service

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/repository"
)

type TaskService interface {
	CreateTask(taskToCreate request.Task) (int64, error)
	GetAllTasks(userId int64, flag int, queryParams request.TaskQueryParams) ([]response.Task, error)
	GetTasksofTeam(teamId int64, queryParams request.TaskQueryParams) ([]response.Task, error)
	UpdateTask(taskToUpdate request.UpdateTask) error
}

type taskService struct {
	taskRepository repository.TaskRepository
}

func NewTaskService(taskRepository repository.TaskRepository) TaskService {
	return taskService{
		taskRepository: taskRepository,
	}
}

func (t taskService) CreateTask(taskToCreate request.Task) (int64, error) {
	return t.taskRepository.CreateTask(taskToCreate)
}

func (t taskService) GetAllTasks(userId int64, flag int, queryParams request.TaskQueryParams) ([]response.Task, error) {
	return t.taskRepository.GetAllTasks(userId, flag, queryParams)
}

func (t taskService) GetTasksofTeam(teamId int64, queryParams request.TaskQueryParams) ([]response.Task, error) {
	return t.taskRepository.GetTasksofTeam(teamId, queryParams)
}

func (t taskService) UpdateTask(taskToUpdate request.UpdateTask) error {
	return t.taskRepository.UpdateTask(taskToUpdate)
}
