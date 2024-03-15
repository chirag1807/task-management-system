package service

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/repository"
)

type TaskService interface{
	CreateTask(taskToCreate request.Task) (int64, error)
	GetAllTasks(userId int64, flag string) ([]response.Task, error)
	GetTasksofTeam(teamId int64) ([]response.Task, error)
	UpdateTask(taskToUpdate request.Task) (error)
	DeleteTask(taskId int64) (error)
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

func (t taskService) GetAllTasks(userId int64, flag string) ([]response.Task, error) {
	return t.taskRepository.GetAllTasks(userId, flag)
}

func (t taskService) GetTasksofTeam(teamId int64) ([]response.Task, error) {
	return t.taskRepository.GetTasksofTeam(teamId)
}

func (t taskService) UpdateTask(taskToUpdate request.Task) (error) {
	return t.taskRepository.UpdateTask(taskToUpdate)
}

func (t taskService) DeleteTask(taskId int64) (error) {
	return t.taskRepository.DeleteTask(taskId)
}