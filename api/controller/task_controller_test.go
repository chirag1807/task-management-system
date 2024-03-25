package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/constant"
	"github.com/go-chi/chi/v5"
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
			StatusCode:         200,
		},
		{
			TestCaseName:       "Field Must be Not Empty.",
			Title:              "",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			Status:             "TO-DO",
			Priority:           "High",
			CreatedBy:          954488202459119617,
			StatusCode:         400,
		},
		{
			TestCaseName:       "Field Must be Required",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			Status:             "TO-DO",
			Priority:           "High",
			CreatedBy:          954488202459119617,
			StatusCode:         400,
		},
		{
			TestCaseName:       "Value Must be in Enum Values.",
			Title:              "Task1",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			Status:             "TO-DO",
			Priority:           "Very Low",
			CreatedBy:          954488202459119617,
			StatusCode:         400,
		},
		{
			TestCaseName: "Assignee Must be User or Team ID(Number).",
			Title:        "Task1",
			Description:  "This is Dummy Task For Test-Cases.",
			Deadline:     time.Now().Add(3 * 24 * time.Hour),
			Status:       "TO-DO",
			Priority:     "Very Low",
			CreatedBy:    954488202459119617,
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/task/create-task", NewTaskController(taskService).CreateTask)

			task := request.Task{
				Title:              v.Title,
				Description:        v.Description,
				Deadline:           v.Deadline,
				AssigneeIndividual: v.AssigneeIndividual,
				AssigneeTeam:       v.AssigneeTeam,
				Status:             v.Status,
				Priority:           v.Priority,
			}
			jsonValue, err := json.Marshal(task)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("POST", "/api/task/create-task", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), constant.UserIdKey, v.CreatedBy)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
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
			StatusCode: 200,
		},
		{
			TestCaseName: "Field Must Be In Enum Values.",
			Flag:         1,
			UserId:       954488202459119617,
			QueryParams: request.TaskQueryParams{
				Limit:        1,
				Offset:       0,
				Search:       "",
				Status:       "Done",
				SortByFilter: true,
			},
			StatusCode: 400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Get("/api/task/get-all-tasks/:Flag", NewTaskController(taskService).GetAllTasks)

			req, err := http.NewRequest("GET", "/api/task/get-all-tasks/:Flag", http.NoBody)
			if err != nil {
				log.Println(err)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("Flag", strconv.Itoa(v.Flag))
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			ctx = context.WithValue(ctx, constant.UserIdKey, v.UserId)
			req = req.WithContext(ctx)

			q := req.URL.Query()
			q.Add("limit", strconv.Itoa(v.QueryParams.Limit))
			q.Add("offset", strconv.Itoa(v.QueryParams.Offset))
			q.Add("search", v.QueryParams.Search)
			q.Add("status", v.QueryParams.Status)
			q.Add("sortByFilter", strconv.FormatBool(v.QueryParams.SortByFilter))
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
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
		{
			TestCaseName: "Field Must Be In Enum Values.",
			TeamID:       954507580144451585,
			QueryParams: request.TaskQueryParams{
				Limit:        1,
				Offset:       0,
				Search:       "",
				Status:       "Done",
				SortByFilter: true,
			},
			StatusCode: 400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Get("/api/task/get-tasks-of-team/:TeamID", NewTaskController(taskService).GetTasksofTeam)

			req, err := http.NewRequest("GET", "/api/task/get-tasks-of-team/:TeamID", http.NoBody)
			if err != nil {
				log.Println(err)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("TeamID", strconv.Itoa(int(v.TeamID)))
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			q := req.URL.Query()
			q.Add("limit", strconv.Itoa(v.QueryParams.Limit))
			q.Add("offset", strconv.Itoa(v.QueryParams.Offset))
			q.Add("search", v.QueryParams.Search)
			q.Add("status", v.QueryParams.Status)
			q.Add("sortByFilter", strconv.FormatBool(v.QueryParams.SortByFilter))
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
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
		Expected           interface{}
		StatusCode         int
	}{
		{
			TestCaseName:       "Task Created Successfully",
			ID:                 954511608047501313,
			Title:              "Task5",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(10 * 24 * time.Hour),
			Status:             "TO-DO",
			Priority:           "High",
			CreatedBy:          954488202459119617,
			StatusCode:         200,
		},
		{
			TestCaseName:       "Field Must be Required",
			Title:              "Task5",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			Status:             "TO-DO",
			Priority:           "High",
			CreatedBy:          954488202459119617,
			StatusCode:         400,
		},
		{
			TestCaseName:       "Value Must be in Enum Values.",
			ID:                 954511608047501313,
			Title:              "Task5",
			Description:        "This is Dummy Task For Test-Cases.",
			Deadline:           time.Now().Add(3 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { id := int64(954497896847212545); return &id }(),
			Status:             "TO-DO",
			Priority:           "Very Low",
			CreatedBy:          954488202459119617,
			StatusCode:         400,
		},
		{
			TestCaseName: "Field Violates Minimun Length Constraint.",
			Title:        "Ta",
			Description:  "This is Dummy Task For Test-Cases.",
			Deadline:     time.Now().Add(3 * 24 * time.Hour),
			Status:       "TO-DO",
			Priority:     "Very Low",
			CreatedBy:    954488202459119617,
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Violates Maximum Length Constraint.",
			Title:        "Task on Backend includes golang, node.js, python etc.",
			Description:  "This is Dummy Task For Test-Cases.",
			Deadline:     time.Now().Add(3 * 24 * time.Hour),
			Status:       "TO-DO",
			Priority:     "Very Low",
			CreatedBy:    954488202459119617,
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Put("/api/task/update-task", NewTaskController(taskService).UpdateTask)

			task := request.Task{
				ID: v.ID,
				Title:              v.Title,
				Description:        v.Description,
				Deadline:           v.Deadline,
				AssigneeIndividual: v.AssigneeIndividual,
				AssigneeTeam:       v.AssigneeTeam,
				Status:             v.Status,
				Priority:           v.Priority,
			}
			jsonValue, err := json.Marshal(task)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("PUT", "/api/task/update-task", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), constant.UserIdKey, v.CreatedBy)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}
