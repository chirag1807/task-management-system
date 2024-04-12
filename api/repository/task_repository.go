package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils/socket"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5"
)

type TaskRepository interface {
	CreateTask(taskToCreate request.Task) (int64, error)
	GetAllTasks(userId int64, queryParams request.TaskQueryParams) ([]response.Task, error)
	//flag is used for get my created tasks and get tasks assigned to me.
	GetTasksofTeam(teamId int64, queryParams request.TaskQueryParams) ([]response.Task, error)
	UpdateTask(taskToUpdate request.UpdateTask) error
}

type taskRepository struct {
	dbConn       *pgx.Conn
	redisClient  *redis.Client
	socketServer *socketio.Server
}

func NewTaskRepo(dbConn *pgx.Conn, redisClient *redis.Client, socketServer *socketio.Server) TaskRepository {
	return taskRepository{
		dbConn:       dbConn,
		redisClient:  redisClient,
		socketServer: socketServer,
	}
}

func (t taskRepository) CreateTask(taskToCreate request.Task) (int64, error) {
	var dbUserPrivacy string
	var dbTeamPrivacy string

	if taskToCreate.AssigneeIndividual != nil {
		rows := t.dbConn.QueryRow(context.Background(), `SELECT privacy FROM users WHERE id = $1`, *taskToCreate.AssigneeIndividual)
		err := rows.Scan(&dbUserPrivacy)
		if err != nil {
			return 0, err
		}
		if dbUserPrivacy != "PUBLIC" {
			return 0, errorhandling.OnlyPublicUserAssignne
		}
	} else if taskToCreate.AssigneeTeam != nil {
		rows := t.dbConn.QueryRow(context.Background(), `SELECT team_privacy FROM teams WHERE id = $1`, taskToCreate.AssigneeTeam)
		err := rows.Scan(&dbTeamPrivacy)
		if err != nil {
			return 0, err
		}
		if dbTeamPrivacy != "PUBLIC" {
			return 0, errorhandling.OnlyPublicTeamAssignne
		}
	}

	var taskId int64
	err := t.dbConn.QueryRow(context.Background(), `INSERT INTO tasks (title, description, deadline, assignee_individual, assignee_team, status, priority,
			created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, taskToCreate.Title, taskToCreate.Description, taskToCreate.Deadline,
		taskToCreate.AssigneeIndividual, taskToCreate.AssigneeTeam, taskToCreate.Status, taskToCreate.Priority, taskToCreate.CreatedBy, taskToCreate.CreatedAt).
		Scan(&taskId)
	if err != nil {
		return 0, err
	}

	taskToCreate.ID = taskId
	taskJSON, _ := json.Marshal(taskToCreate)
	err = t.redisClient.Set(context.Background(), "tasks:"+strconv.FormatInt(taskId, 10), taskJSON, 0).Err()
	if err != nil {
		return 0, err
	}

	if taskToCreate.AssigneeTeam != nil {
		err = t.redisClient.SAdd(context.Background(), "tasks:assigned_to_team:"+strconv.FormatInt(*taskToCreate.AssigneeTeam, 10), taskId).Err()
		if err != nil {
			return 0, err
		}
	}
	if taskToCreate.AssigneeIndividual != nil {
		err = t.redisClient.SAdd(context.Background(), "tasks:assigned_to_user:"+strconv.FormatInt(*taskToCreate.AssigneeIndividual, 10), taskId).Err()
		if err != nil {
			return 0, err
		}
	}

	err = t.redisClient.SAdd(context.Background(), "tasks:created_by:"+strconv.FormatInt(taskToCreate.CreatedBy, 10), taskId).Err()
	if err != nil {
		return 0, err
	}

	if taskToCreate.AssigneeIndividual != nil {
		socket.EmitCreateAndUpdateTaskEvents(t.socketServer, strconv.FormatInt(*taskToCreate.AssigneeIndividual, 10), constant.EMPTY_STRING, taskToCreate, 0)
	}
	if taskToCreate.AssigneeTeam != nil {
		socket.EmitCreateAndUpdateTaskEvents(t.socketServer, "task-created", strconv.FormatInt(*taskToCreate.AssigneeTeam, 10), taskToCreate, 1)
	}

	return taskId, nil
}

func (t taskRepository) GetAllTasks(userId int64, queryParams request.TaskQueryParams) ([]response.Task, error) {
	var tasks pgx.Rows
	var err error
	tasksSlice := make([]response.Task, 0)
	var query string

	if !queryParams.CreatedByMe {
		taskIds, _ := t.redisClient.SMembers(context.Background(), "tasks:created_by:"+strconv.FormatInt(userId, 10)).Result()
		tasksSlice, _ = GetTasksFromRedisByIDList(t.redisClient, taskIds)
		if len(tasksSlice) != 0 {
			return tasksSlice, nil
		} else {
			query = `SELECT * FROM tasks WHERE created_by = $1 AND true`
			query = CreateQueryForParamsOfGetTask(query, queryParams)
			tasks, err = t.dbConn.Query(context.Background(), query, userId)
			if err != nil {
				return tasksSlice, err
			}
		}
	}
	if queryParams.CreatedByMe {
		taskIds, _ := t.redisClient.SMembers(context.Background(), "tasks:assigned_to_user:"+strconv.FormatInt(userId, 10)).Result()
		userTeams, _ := t.redisClient.SMembers(context.Background(), "user:"+strconv.FormatInt(userId, 10)+":teams").Result()
		for _, teamId := range userTeams {
			teamTaskIDs, _ := t.redisClient.SMembers(context.Background(), "tasks:assigned_to_team:"+teamId).Result()
			taskIds = append(taskIds, teamTaskIDs...)
		}
		tasksSlice, err = GetTasksFromRedisByIDList(t.redisClient, taskIds)
		if err != nil {
			return nil, err
		}
		if len(tasksSlice) != 0 {
			return tasksSlice, nil
		} else {
			query = `SELECT * FROM tasks WHERE (assignee_individual = $1 OR assignee_team IN (SELECT team_id from team_members where member_id = $2))`
			query = CreateQueryForParamsOfGetTask(query, queryParams)
			tasks, err = t.dbConn.Query(context.Background(), query, userId, userId)
			if err != nil {
				return tasksSlice, err
			}
		}
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
	teamTaskIDs, _ := t.redisClient.SMembers(context.Background(), "tasks:assigned_to_team:"+strconv.FormatInt(teamId, 10)).Result()
	tasksSlice, _ := GetTasksFromRedisByIDList(t.redisClient, teamTaskIDs)
	if len(tasksSlice) != 0 {
		return tasksSlice, nil
	}

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
	if queryParams.Search != constant.EMPTY_STRING {
		query += fmt.Sprintf(" AND (title ILIKE '%%%s%%' OR description ILIKE '%%%s%%')", queryParams.Search, queryParams.Search)
	}
	if queryParams.Status != constant.EMPTY_STRING {
		query += fmt.Sprintf(" AND status = '%s'", queryParams.Status)
	}
	if queryParams.SortByFilter {
		query += " ORDER BY CASE priority WHEN 'VERY HIGH' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'MEDIUM' THEN 3 ELSE 4 END"
	}
	query += fmt.Sprintf(" LIMIT %d", queryParams.Limit)
	query += fmt.Sprintf(" OFFSET %d", queryParams.Offset)
	return query
}

func (t taskRepository) UpdateTask(taskToUpdate request.UpdateTask) error {
	var dbTask response.Task
	rows := t.dbConn.QueryRow(context.Background(), `SELECT * FROM tasks WHERE id = $1`, taskToUpdate.ID)

	err := rows.Scan(&dbTask.ID, &dbTask.Title, &dbTask.Description, &dbTask.Deadline, &dbTask.AssigneeIndividual, &dbTask.AssigneeTeam, &dbTask.Status,
		&dbTask.Priority, &dbTask.CreatedBy, &dbTask.CreatedAt, &dbTask.UpdatedBy, &dbTask.UpdatedAt)

	if err != nil {
		if err.Error() == constant.PG_NO_ROWS {
			return errorhandling.NoTaskFound
		}
		return err
	}

	if dbTask.Status == "Closed" {
		return errorhandling.TaskClosed
	}

	if dbTask.AssigneeIndividual != nil {
		if *dbTask.AssigneeIndividual != *taskToUpdate.UpdatedBy && dbTask.CreatedBy != *taskToUpdate.UpdatedBy {
			return errorhandling.NotAllowed
		}
	} else {
		var userCount int
		rows := t.dbConn.QueryRow(context.Background(), `SELECT COUNT(*) FROM team_members WHERE team_id = $1 AND member_id = $2`, dbTask.AssigneeTeam, taskToUpdate.UpdatedBy)
		err := rows.Scan(&userCount)
		if err != nil {
			return err
		}
		if userCount == 0 && dbTask.CreatedBy != *taskToUpdate.UpdatedBy {
			return errorhandling.NotAllowed
		}
	}

	if dbTask.CreatedBy != *taskToUpdate.UpdatedBy && (taskToUpdate.Priority != constant.EMPTY_STRING || taskToUpdate.Title != constant.EMPTY_STRING || taskToUpdate.Description != constant.EMPTY_STRING ||
		taskToUpdate.AssigneeIndividual != nil || taskToUpdate.AssigneeTeam != nil) {
		return errorhandling.NotAllowed
	}

	query, args, err := UpdateQuery("tasks", taskToUpdate, taskToUpdate.ID, 1)
	if err != nil {
		return err
	}
	_, err = t.dbConn.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	taskToUpdateinRedis := UpdateTaskFields(dbTask, taskToUpdate)
	taskJSON, err := json.Marshal(taskToUpdateinRedis)
	if err != nil {
		return err
	}
	t.redisClient.Set(context.Background(), "tasks:"+strconv.FormatInt(taskToUpdateinRedis.ID, 10), taskJSON, 0)

	if taskToUpdate.AssigneeIndividual != nil || taskToUpdate.AssigneeTeam != nil {
		if dbTask.AssigneeTeam != nil {
			if taskToUpdate.AssigneeTeam != nil {
				t.redisClient.SRem(context.Background(), "tasks:assigned_to_team:"+strconv.FormatInt(*dbTask.AssigneeTeam, 10), taskToUpdate.ID)
				t.redisClient.SAdd(context.Background(), "tasks:assigned_to_team:"+strconv.FormatInt(*taskToUpdate.AssigneeTeam, 10), taskToUpdate.ID)
			} else {
				t.redisClient.SRem(context.Background(), "tasks:assigned_to_team:"+strconv.FormatInt(*dbTask.AssigneeTeam, 10), taskToUpdate.ID)
				t.redisClient.SAdd(context.Background(), "tasks:assigned_to_user:"+strconv.FormatInt(*taskToUpdate.AssigneeIndividual, 10), taskToUpdate.ID)
			}
		}
		if dbTask.AssigneeIndividual != nil {
			if taskToUpdate.AssigneeTeam != nil {
				t.redisClient.SRem(context.Background(), "tasks:assigned_to_user:"+strconv.FormatInt(*dbTask.AssigneeIndividual, 10), taskToUpdate.ID)
				t.redisClient.SAdd(context.Background(), "tasks:assigned_to_team:"+strconv.FormatInt(*taskToUpdate.AssigneeTeam, 10), taskToUpdate.ID)
			} else {
				t.redisClient.SRem(context.Background(), "tasks:assigned_to_team:"+strconv.FormatInt(*dbTask.AssigneeIndividual, 10), taskToUpdate.ID)
				t.redisClient.SAdd(context.Background(), "tasks:assigned_to_user:"+strconv.FormatInt(*taskToUpdate.AssigneeIndividual, 10), taskToUpdate.ID)
			}
		}
	}

	if dbTask.AssigneeIndividual != nil {
		if taskToUpdate.AssigneeIndividual != nil {
			socket.EmitCreateAndUpdateTaskEvents(t.socketServer, strconv.FormatInt(*taskToUpdate.AssigneeIndividual, 10), constant.EMPTY_STRING, taskToUpdateinRedis, 0)
		} else if taskToUpdate.AssigneeTeam != nil {
			socket.EmitCreateAndUpdateTaskEvents(t.socketServer, "task-updated", strconv.FormatInt(*taskToUpdate.AssigneeTeam, 10), taskToUpdateinRedis, 1)
		} else {
			socket.EmitCreateAndUpdateTaskEvents(t.socketServer, strconv.FormatInt(*dbTask.AssigneeIndividual, 10), constant.EMPTY_STRING, taskToUpdateinRedis, 0)
		}
	}
	if dbTask.AssigneeTeam != nil {
		if taskToUpdate.AssigneeIndividual != nil {
			socket.EmitCreateAndUpdateTaskEvents(t.socketServer, strconv.FormatInt(*taskToUpdate.AssigneeIndividual, 10), constant.EMPTY_STRING, taskToUpdateinRedis, 0)
		} else if taskToUpdate.AssigneeTeam != nil {
			socket.EmitCreateAndUpdateTaskEvents(t.socketServer, "task-updated", strconv.FormatInt(*taskToUpdate.AssigneeTeam, 10), taskToUpdateinRedis, 1)
		} else {
			socket.EmitCreateAndUpdateTaskEvents(t.socketServer, "task-updated", strconv.FormatInt(*dbTask.AssigneeTeam, 10), taskToUpdateinRedis, 1)
		}
	}

	return nil
}

func GetTasksFromRedisByIDList(redisClient *redis.Client, taskIds []string) ([]response.Task, error) {
	var tasks []response.Task
	for _, taskId := range taskIds {
		taskJSON, err := redisClient.Get(context.Background(), "tasks:"+taskId).Result()
		if err != nil {
			return tasks, err
		}

		var task response.Task
		json.Unmarshal([]byte(taskJSON), &task)

		tasks = append(tasks, task)
	}

	return tasks, nil
}
