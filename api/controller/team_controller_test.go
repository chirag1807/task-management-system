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

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/constant"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Name         string
		CreatedBy    int64
		Privacy  *string
		MemberIDs     []int64
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Created Successfully",
			Name:         "Team Mahakal",
			MemberIDs:     []int64{954497896847212545},
			CreatedBy:    954488202459119617,
			StatusCode:   200,
		},
		{
			TestCaseName: "Field Must be Not Empty.",
			Name:         "",
			CreatedBy:    954488202459119617,
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Must be Required",
			CreatedBy:    954488202459119617,
			StatusCode:   400,
		},
		{
			TestCaseName: "Value Must be in Enum Values.",
			Name:         "Team Mahakal",
			Privacy:  func() *string { team_privacy := string("Protected"); return &team_privacy }(),
			CreatedBy:    954488202459119617,
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/v1/teams", NewTeamController(teamService).CreateTeam)

			task := request.CreateTeam{
				Details: request.Team{
					Name:        v.Name,
					Privacy: v.Privacy,
				},
				Members: v.MemberIDs,
			}
			jsonValue, err := json.Marshal(task)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("POST", "/api/v1/teams", bytes.NewBuffer(jsonValue))
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

func TestAddMembersToTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName     string
		TeamID           int64
		MemberID         []int64
		GroupCreatedByID int64
		Expected         interface{}
		StatusCode       int
	}{
		{
			TestCaseName:     "Team Members Added Successfully",
			TeamID:           954507580144451585,
			MemberID:         []int64{954497896847212545},
			GroupCreatedByID: 954488202459119617,
			StatusCode:       200,
		},
		{
			TestCaseName:     "Field Must be Required",
			MemberID:         []int64{954497896847212545},
			GroupCreatedByID: 954488202459119617,
			StatusCode:       400,
		},
		{
			TestCaseName:     "Member Already Exist",
			TeamID:           954507580144451585,
			MemberID:         []int64{954497896847212545},
			GroupCreatedByID: 954488202459119617,
			StatusCode:       409,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/v1/teams/:TeamID/members", NewTeamController(teamService).AddMembersToTeam)

			task := request.TeamMembersWithTeamID{
				TeamID:   v.TeamID,
				MemberIDs: v.MemberID,
			}
			jsonValue, err := json.Marshal(task)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("POST", "/api/v1/teams/:TeamID/members", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("TeamID", strconv.Itoa(int(v.TeamID)))
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			ctx = context.WithValue(ctx, constant.UserIdKey, v.GroupCreatedByID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestRemoveMembersFromTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName     string
		TeamID           int64
		MemberID         []int64
		GroupCreatedByID int64
		Expected         interface{}
		StatusCode       int
	}{
		{
			TestCaseName:     "Team Members Removed Successfully",
			TeamID:           954507580144451585,
			MemberID:         []int64{954497896847212545},
			GroupCreatedByID: 954488202459119617,
			StatusCode:       200,
		},
		{
			TestCaseName:     "Field Must be Required",
			MemberID:         []int64{954497896847212545},
			GroupCreatedByID: 954488202459119617,
			StatusCode:       400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Delete("/api/v1/teams/:TeamID/members", NewTeamController(teamService).RemoveMembersFromTeam)

			task := request.TeamMembersWithTeamID{
				TeamID:   v.TeamID,
				MemberIDs: v.MemberID,
			}
			jsonValue, err := json.Marshal(task)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("DELETE", "/api/v1/teams/:TeamID/members", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("TeamID", strconv.Itoa(int(v.TeamID)))
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			ctx = context.WithValue(ctx, constant.UserIdKey, v.GroupCreatedByID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestGetAllTeams(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		UserId       int64
		QueryParams  request.TeamQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Created By Me - Success",
			UserId:       954488202459119617,
			QueryParams: request.TeamQueryParams{
				CreatedByMe: true,
				Limit:          1,
				Offset:         0,
				Search:         "",
				SortByCreatedAt: true,
			},
			StatusCode: 200,
		},
		{
			TestCaseName: "Team in Which I were Added - Success",
			UserId:       954488202459119617,
			QueryParams: request.TeamQueryParams{
				CreatedByMe: false,
				Limit:          1,
				Offset:         0,
				Search:         "",
				SortByCreatedAt: true,
			},
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Get("/api/v1/teams", NewTeamController(teamService).GetAllTeams)

			req, err := http.NewRequest("GET", "/api/v1/teams", http.NoBody)
			if err != nil {
				log.Println(err)
			}

			ctx := context.WithValue(req.Context(), constant.UserIdKey, v.UserId)
			req = req.WithContext(ctx)

			q := req.URL.Query()
			q.Add("createdByMe", strconv.FormatBool(v.QueryParams.CreatedByMe))
			q.Add("limit", strconv.Itoa(v.QueryParams.Limit))
			q.Add("offset", strconv.Itoa(v.QueryParams.Offset))
			q.Add("search", v.QueryParams.Search)
			q.Add("sortByCreatedAt", strconv.FormatBool(v.QueryParams.SortByCreatedAt))
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestGetTeamMembers(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		TeamID       int64
		QueryParams  request.TeamQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Members Fetched Successfully",
			TeamID:       954507580144451585,
			QueryParams: request.TeamQueryParams{
				Limit:  1,
				Offset: 0,
			},
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Get("/api/v1/teams/:TeamID/members", NewTeamController(teamService).GetTeamMembers)

			req, err := http.NewRequest("GET", "/api/v1/teams/:TeamID/members", http.NoBody)
			if err != nil {
				log.Println(err)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("TeamID", strconv.FormatInt(v.TeamID, 10))
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			q := req.URL.Query()
			q.Add("limit", strconv.Itoa(v.QueryParams.Limit))
			q.Add("offset", strconv.Itoa(v.QueryParams.Offset))
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestLeftTeam(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		TeamID       int64
		UserID       int64
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Team Left Successfully",
			TeamID:       954507580144451585,
			UserID:       954488202459119617,
			StatusCode:   200,
		},
		{
			TestCaseName: "Can not Left Team Because Not a Member",
			TeamID:       954507580144451585,
			UserID:       954497896847212545,
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Get("/api/v1/teams/leave/:TeamID", NewTeamController(teamService).LeaveTeam)

			req, err := http.NewRequest("GET", "/api/v1/teams/leave/:TeamID", http.NoBody)
			if err != nil {
				log.Println(err)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("TeamID", strconv.FormatInt(v.TeamID, 10))
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			ctx = context.WithValue(ctx, constant.UserIdKey, v.UserID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}
