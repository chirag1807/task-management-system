package repository

import (
	"context"
	"fmt"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type TaskRepository interface {
	CreateTask(taskToCreate request.Task) (int64, error)
	GetAllTasks(userId int64, flag int, queryParams request.TaskQueryParams) ([]response.Task, error)
	//flag is used for get my created tasks and get tasks assigned to me.
	GetTasksofTeam(teamId int64, queryParams request.TaskQueryParams) ([]response.Task, error)
	UpdateTask(taskToUpdate request.Task) error
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

func (t taskRepository) GetAllTasks(userId int64, flag int, queryParams request.TaskQueryParams) ([]response.Task, error) {
	var tasks pgx.Rows
	var err error
	tasksSlice := make([]response.Task, 0)

	var query string
	if flag == 0 {
		query = `SELECT * FROM tasks WHERE created_by = $1 AND true`
		query = CreateQueryForParamsOfGetTask(query, queryParams)
		tasks, err = t.dbConn.Query(context.Background(), query, userId)
	}
	if flag == 1 {
		query = `SELECT * FROM tasks WHERE (assignee_individual = $1 OR assignee_team IN (SELECT team_id from team_members where member_id = $2))`
		query = CreateQueryForParamsOfGetTask(query, queryParams)
		tasks, err = t.dbConn.Query(context.Background(), query, userId, userId)
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

func (t taskRepository) GetTasksofTeam(teamId int64, queryParams request.TaskQueryParams) ([]response.Task, error) {
	tasksSlice := make([]response.Task, 0)

	query := `SELECT * FROM tasks WHERE assignee_team = $1 AND true`
	query = CreateQueryForParamsOfGetTask(query, queryParams)
	tasks, err := t.dbConn.Query(context.Background(), query, teamId)
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

func CreateQueryForParamsOfGetTask(query string, queryParams request.TaskQueryParams) string {
	if queryParams.Search != "" {
		query += fmt.Sprintf(" AND (title ILIKE '%%%s%%' OR description ILIKE '%%%s%%')", queryParams.Search, queryParams.Search)
	}
	if queryParams.Status != "" {
		query += fmt.Sprintf(" AND status = '%s'", queryParams.Status)
	}
	if queryParams.SortByFilter {
		query += " ORDER BY CASE priority WHEN 'Very High' THEN 1 WHEN 'High' THEN 2 WHEN 'Medium' THEN 3 ELSE 4 END"
	}
	query += fmt.Sprintf(" LIMIT %d", queryParams.Limit)
	query += fmt.Sprintf(" OFFSET %d", queryParams.Offset)
	return query
}

func (t taskRepository) UpdateTask(taskToUpdate request.Task) error {
	var dbTask response.Task
	task := t.dbConn.QueryRow(context.Background(), `SELECT * FROM tasks WHERE id = $1`, taskToUpdate.ID)
	err := task.Scan(&dbTask.ID, &dbTask.Title, &dbTask.Description, &dbTask.Deadline, &dbTask.AssigneeIndividual, &dbTask.AssigneeTeam, &dbTask.Status, 
		&dbTask.Priority, &dbTask.CreatedBy, &dbTask.CreatedAt, &dbTask.UpdatedBy, &dbTask.UpdatedAt)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return errorhandling.NoTaskFound
		}
		return err
	}

	if dbTask.Status == "Closed" {
		return errorhandling.TaskClosed
	}

	if dbTask.AssigneeIndividual != nil {
		if *dbTask.AssigneeIndividual != *taskToUpdate.UpdatedBy {
			return errorhandling.NotAllowed
		}
	} else {
		var userCount int
		t.dbConn.QueryRow(context.Background(), `SELECT COUNT(*) FROM team_members WHERE team_id = $1 AND member_id = $2`, dbTask.AssigneeTeam, taskToUpdate.UpdatedBy).Scan(&userCount)
		if userCount == 0 {
			return errorhandling.NotAllowed
		}
	}

	if dbTask.CreatedBy != *taskToUpdate.UpdatedBy && (taskToUpdate.Priority != "" || taskToUpdate.Title != "" || taskToUpdate.Description != "" || 
		taskToUpdate.AssigneeIndividual != nil || taskToUpdate.AssigneeTeam != nil) {
		return errorhandling.NotAllowed
	}

	query, args, err := UpdateQuery("tasks", taskToUpdate, taskToUpdate.ID)
	if err != nil {
		return err
	}
	_, err = t.dbConn.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	return nil
}
