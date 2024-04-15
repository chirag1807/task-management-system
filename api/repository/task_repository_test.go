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
			TestCaseName:       "Task Created for Assignee Individual Successfully",
			Title:              "Task1",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			Status:             "TO-DO",
			Priority:           "HIGH",
			CreatedBy:          954488202459119617,
			Expected:           nil,
			StatusCode:         200,
		},
		{
			TestCaseName: "Task Created for Assignee Team Successfully",
			Title:        "Task1",
			Description:  "This is Dummy Task For Test-Cases.",
			Deadline:     time.Now().Add(3 * 24 * time.Hour),
			AssigneeTeam: func() *int64 { id := int64(954507580144451585); return &id }(),
			Status:       "TO-DO",
			Priority:     "HIGH",
			CreatedBy:    954488202459119617,
			Expected:     nil,
			StatusCode:   200,
		},
		{
			TestCaseName:       "Task Can be Assigned to Public User Only",
			Title:              "Task1",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212546); return &id }(),
			Status:             "TO-DO",
			Priority:           "HIGH",
			CreatedBy:          954488202459119617,
			Expected:           errorhandling.OnlyPublicUserAssignne,
			StatusCode:         400,
		},
		{
			TestCaseName: "Task Can be Assigned to Public Team Only",
			Title:        "Task1",
			Description:  "This is Dummy Task For Test-Cases.",
			Deadline:     time.Now().Add(3 * 24 * time.Hour),
			AssigneeTeam: func() *int64 { id := int64(954507580144451586); return &id }(),
			Status:       "TO-DO",
			Priority:     "HIGH",
			CreatedBy:    954488202459119617,
			Expected:     errorhandling.OnlyPublicTeamAssignne,
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			task := request.Task{
				Title:              v.Title,
				Description:        v.Description,
				Deadline:           v.Deadline,
				AssigneeIndividual: v.AssigneeIndividual,
				AssigneeTeam:       v.AssigneeTeam,
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
		UserId       int64
		QueryParams  request.TaskQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Task Created By Me - Success",
			UserId:       954488202459119617,
			QueryParams: request.TaskQueryParams{
				CreatedByMe:  false,
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
			UserId:       954488202459119617,
			QueryParams: request.TaskQueryParams{
				CreatedByMe:  true,
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

			_, err := NewTaskRepo(dbConn, redisClient, socketServer).GetAllTasks(v.UserId, v.QueryParams)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestCreateQueryForParamsOfGetTask(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		QueryParams  request.TaskQueryParams
		Expected     interface{}
	}{
		{
			TestCaseName: "Query Based on Query Params Created.",
			QueryParams: request.TaskQueryParams{
				Limit:        10,
				Offset:       0,
				Search:       "Chirag",
				Status:       "TO-DO",
				SortByFilter: true,
			},
			Expected: ` AND (title ILIKE '%Chirag%' OR description ILIKE '%Chirag%') AND status = 'TO-DO' ORDER BY CASE priority WHEN 'VERY HIGH' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'MEDIUM' THEN 3 ELSE 4 END LIMIT 10 OFFSET 0`,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			query := CreateQueryForParamsOfGetTask("", v.QueryParams)
			assert.Equal(t, v.Expected, query)
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
			TestCaseName: "Task Updated Successfully",
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
		{
			TestCaseName: "Task is Closed",
			ID:           954511608047501314,
			Title:        "Task123",
			CreatedBy:    954488202459119617,
			UpdatedBy:    954488202459119617,
			UpdatedAt:    time.Now(),
			Expected:     errorhandling.TaskClosed,
			StatusCode:   400,
		},
		{
			TestCaseName:       "Assignee Team to Assignee Individual",
			ID:                 954511608047501313,
			Title:              "Task123",
			CreatedBy:          954488202459119617,
			UpdatedBy:          954488202459119617,
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			UpdatedAt:          time.Now(),
			Expected:           nil,
			StatusCode:         200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			task := request.UpdateTask{
				ID:                 v.ID,
				Title:              v.Title,
				AssigneeIndividual: v.AssigneeIndividual,
				UpdatedBy:          &v.UpdatedBy,
				UpdatedAt:          &v.UpdatedAt,
			}

			err := NewTaskRepo(dbConn, redisClient, socketServer).UpdateTask(task)
			fmt.Println(err)
			assert.Equal(t, v.Expected, err)
		})
	}
}
