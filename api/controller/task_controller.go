package controller

import (
	"net/http"

	"github.com/chirag1807/task-management-system/api/service"
)

type TaskController interface{
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetAllTasks(w http.ResponseWriter, r *http.Request)
	GetTasksofTeam(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
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

}

func (t taskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	
}

func (t taskController) GetTasksofTeam(w http.ResponseWriter, r *http.Request) {
	
}

func (t taskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	
}

func (t taskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	
}