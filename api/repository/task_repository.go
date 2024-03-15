package repository

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type TaskRepository interface{
	CreateTask(taskToCreate request.Task) (int64, error)
	GetAllTasks(userId int64, flag string) ([]response.Task, error)
	//flag is used for get my created tasks and get tasks assigned to me.
	GetTasksofTeam(teamId int64) ([]response.Task, error)
	UpdateTask(taskToUpdate request.Task) (error)
	DeleteTask(taskId int64) (error)
}

type taskRepository struct {
	dbConn      *pgx.Conn
	redisClient *redis.Client
}

func NewTaskRepo(dbConn *pgx.Conn, redisClient *redis.Client) TaskRepository {
	return taskRepository{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (t taskRepository) CreateTask(taskToCreate request.Task) (int64, error) {
	return 1, nil
}

func (t taskRepository) GetAllTasks(userId int64, flag string) ([]response.Task, error) {
	return []response.Task{}, nil
}

func (t taskRepository) GetTasksofTeam(teamId int64) ([]response.Task, error) {
	return []response.Task{}, nil
}

func (t taskRepository) UpdateTask(taskToUpdate request.Task) (error) {
	return nil
}

func (t taskRepository) DeleteTask(taskId int64) (error) {
	return nil
}