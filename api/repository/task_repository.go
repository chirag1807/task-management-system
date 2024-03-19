package repository

import (
	"context"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type TaskRepository interface {
	CreateTask(taskToCreate request.Task) (int64, error)
	GetAllTasks(userId int64, flag int) ([]response.Task, error)
	//flag is used for get my created tasks and get tasks assigned to me.
	GetTasksofTeam(teamId int64) ([]response.Task, error)
	UpdateTask(taskToUpdate request.Task) error
	DeleteTask(taskId int64) error
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
	var dbUserprofile string
	if taskToCreate.AssigneeIndividual != nil {
		t.dbConn.QueryRow(context.Background(), `SELECT profile FROM users WHERE id = $1`, taskToCreate.AssigneeIndividual).Scan(&dbUserprofile)
		if dbUserprofile != "Public" {
			return 0, errorhandling.OnlyPublicUserAssignne
		}
	}
	var dbTeamprofile string
	if taskToCreate.AssigneeTeam != nil {
		t.dbConn.QueryRow(context.Background(), `SELECT team_profile FROM teams WHERE id = $1`, taskToCreate.AssigneeTeam).Scan(&dbTeamprofile)
		if dbTeamprofile != "Public" {
			return 0, errorhandling.OnlyPublicTeamAssignne
		}
	}

	var taskID int64
	err := t.dbConn.QueryRow(context.Background(), `INSERT INTO tasks (title, description, deadline, assignee_individual, assignee_team, status, priority,
		created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, taskToCreate.Title, taskToCreate.Description, taskToCreate.Deadline,
		taskToCreate.AssigneeIndividual, taskToCreate.AssigneeTeam, taskToCreate.Status, taskToCreate.Priority, taskToCreate.CreatedBy, taskToCreate.CreatedAt).
		Scan(&taskID)
	if err != nil {
		return 0, err
	}
	return taskID, nil
}

func (t taskRepository) GetAllTasks(userId int64, flag int) ([]response.Task, error) {
	var tasks pgx.Rows
	var err error
	tasksSlice := make([]response.Task, 0)

	if flag == 0 {
		tasks, err = t.dbConn.Query(context.Background(), `SELECT * FROM tasks WHERE created_by = $1`, userId)
	}
	if flag == 1 {
		tasks, err = t.dbConn.Query(context.Background(), `SELECT * FROM tasks WHERE assignee_individual = $1 OR assignee_team
		IN (SELECT team_id from team_members where member_id = $2)`, userId, userId)
	}

	if err != nil {
		return tasksSlice, err
	}
	defer tasks.Close()

	var task response.Task
	for tasks.Next() {
		if err := tasks.Scan(&task.ID, &task.Title, &task.Description, &task.Deadline, &task.AssigneeIndividual, &task.AssigneeTeam, &task.Status,
			&task.Priority, &task.CreatedBy, &task.CreatedAt, &task.UpdatedBy, &task.UpdatedAt); err != nil {
			return tasksSlice, err
		}
		tasksSlice = append(tasksSlice, task)
	}

	return tasksSlice, nil
}

func (t taskRepository) GetTasksofTeam(teamId int64) ([]response.Task, error) {
	tasksSlice := make([]response.Task, 0)

	tasks, err := t.dbConn.Query(context.Background(), `SELECT * FROM tasks WHERE assignee_team = $1`, teamId)
	if err != nil {
		return tasksSlice, err
	}
	defer tasks.Close()

	var task response.Task
	for tasks.Next() {
		if err := tasks.Scan(&task.ID, &task.Title, &task.Description, &task.Deadline, &task.AssigneeIndividual, &task.AssigneeTeam, &task.Status,
			&task.Priority, &task.CreatedBy, &task.CreatedAt, &task.UpdatedBy, &task.UpdatedAt); err != nil {
			return tasksSlice, err
		}
		tasksSlice = append(tasksSlice, task)
	}

	return tasksSlice, nil
}

func (t taskRepository) UpdateTask(taskToUpdate request.Task) error {
	return nil
}

func (t taskRepository) DeleteTask(taskId int64) error {
	return nil
}
