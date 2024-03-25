package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	testCases := []struct {
		TestCaseName       string
		Title              string
		Description        string
		Deadline           time.Time
		AssigneeIndividual *int64
		AssigneeTeam       *int64
		Status             string
		Priority           string
		CreatedBy          int64
		Expected           interface{}
		StatusCode         int
	}{
		{
			TestCaseName:       "Task Created Successfully",
			Title:              "Task1",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			Status:             "TO-DO",
			Priority:           "High",
			CreatedBy:          954488202459119617,
			Expected:           nil,
			StatusCode:         200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			task := request.Task{
				Title:              v.Title,
				Description:        v.Description,
				Deadline:           v.Deadline,
				AssigneeIndividual: v.AssigneeIndividual,
				Status:             v.Status,
				Priority:           v.Priority,
				CreatedBy:          v.CreatedBy,
				CreatedAt:          time.Now(),
			}

			_, err := NewTaskRepo(dbConn, redisClient, socketServer).CreateTask(task)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Flag         int
		UserId       int64
		QueryParams  request.TaskQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Task Created By Me - Success",
			Flag:         0,
			UserId:       954488202459119617,
			QueryParams: request.TaskQueryParams{
				Limit:        1,
				Offset:       0,
				Search:       "",
				Status:       "",
				SortByFilter: true,
			},
			Expected:   nil,
			StatusCode: 200,
		},
		{
			TestCaseName: "Task Assigned To Me - Success",
			Flag:         1,
			UserId:       954488202459119617,
			QueryParams: request.TaskQueryParams{
				Limit:        1,
				Offset:       0,
				Search:       "",
				Status:       "",
				SortByFilter: true,
			},
			Expected:   nil,
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			_, err := NewTaskRepo(dbConn, redisClient, socketServer).GetAllTasks(v.UserId, v.Flag, v.QueryParams)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestGetTasksofTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		TeamID       int64
		QueryParams  request.TaskQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Tasks Of Team - Success",
			TeamID:       954507580144451585,
			QueryParams: request.TaskQueryParams{
				Limit:        1,
				Offset:       0,
				Search:       "",
				Status:       "",
				SortByFilter: true,
			},
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			_, err := NewTaskRepo(dbConn, redisClient, socketServer).GetTasksofTeam(v.TeamID, v.QueryParams)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	testCases := []struct {
		TestCaseName       string
		ID                 int64
		Title              string
		Description        string
		Deadline           time.Time
		AssigneeIndividual *int64
		AssigneeTeam       *int64
		Status             string
		Priority           string
		CreatedBy          int64
		UpdatedBy          int64
		UpdatedAt          time.Time
		Expected           interface{}
		StatusCode         int
	}{
		{
			TestCaseName: "Task Created Successfully",
			ID:           954511608047501313,
			Title:        "Task123",
			CreatedBy:    954488202459119617,
			UpdatedBy:    954488202459119617,
			UpdatedAt:    time.Now(),
			Expected:     nil,
			StatusCode:   200,
		},
		{
			TestCaseName: "Not Allowed to Update Task",
			ID:           954511608047501313,
			Title:        "Task123",
			CreatedBy:    954488202459119617,
			UpdatedBy:    954497896847212545,
			UpdatedAt:    time.Now(),
			Expected:     errorhandling.NotAllowed,
			StatusCode:   403,
		},
		{
			TestCaseName: "No Task Found",
			ID:           1,
			Title:        "Task123",
			CreatedBy:    954488202459119617,
			UpdatedBy:    954488202459119617,
			UpdatedAt:    time.Now(),
			Expected:     errorhandling.NoTaskFound,
			StatusCode:   404,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			task := request.Task{
				ID:        v.ID,
				Title:     v.Title,
				CreatedBy: v.CreatedBy,
				UpdatedBy: &v.UpdatedBy,
				UpdatedAt: &v.UpdatedAt,
			}

			err := NewTaskRepo(dbConn, redisClient, socketServer).UpdateTask(task)
			fmt.Println(err)
			assert.Equal(t, v.Expected, err)
		})
	}
}
