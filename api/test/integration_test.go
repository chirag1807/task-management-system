package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chirag1807/logease"
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/route"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/constant"
	"github.com/chirag1807/task-management-system/db"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/chirag1807/task-management-system/utils/socket"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

var dbConn *pgx.Conn
var redisClient *redis.Client
var rabbitmqConn *amqp.Connection
var socketServer *socketio.Server

func runTestServer() (*httptest.Server, *pgx.Conn) {
	config.LoadConfig("../../.config/", "../../.config/secret.json")
	dbConn, redisClient, rabbitmqConn = db.SetDBConection(1)
	socketServer = socket.SocketConnection()
	loggerInstance, err := logease.InitLogease(false, config.Config.TeamsWebHookURL, logease.Slog)
	if err != nil {
		log.Fatal(err)
	}
	_ = loggerInstance.(logease.SlogLoggerInstance)

	r := route.InitializeRouter(dbConn, redisClient, rabbitmqConn, socketServer)

	err = utils.ClearMockData(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := dbConn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		log.Fatal(err)
	}

	tx, err = utils.InsertMockData(tx)
	if err != nil {
		tx.Rollback(context.Background())
		log.Fatal(err)
	}
	tx.Commit(context.Background())

	return httptest.NewServer(r), dbConn
}

func TestTaskCRUD(t *testing.T) {
	buf := new(bytes.Buffer)
	var id int64 = 0
	ts, dbConn := runTestServer()
	defer ts.Close()

	t.Run("it should return ok when task creation done successfully.", func(t *testing.T) {
		topic := request.Task{
			Title:              "Dummy Task",
			Description:        "This id Dummy Task.",
			Deadline:           time.Now().Add(5 * 24 * time.Hour),
			AssigneeIndividual: func() *int64 { userId := int64(954488202459119617); return &userId }(),
			Status:             "TO-DO",
			Priority:           "HIGH",
		}
		jsonValue, err := json.Marshal(topic)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tasks", ts.URL), bytes.NewBuffer(jsonValue))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		var response response.SuccessResponse
		json.Unmarshal([]byte(responseBody), &response)
		fmt.Println(responseBody)
		id = *response.ID
		assert.Equal(t, constant.TASK_CREATED, response.Message)
	})

	t.Run("it should return tasks assigned to me.", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tasks?createdByMe=true", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		buf.ReadFrom(resp.Body)
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		fmt.Println(responseBody)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("it should return tasks created by me.", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tasks?createdByMe=false", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		buf.ReadFrom(resp.Body)
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		fmt.Println(responseBody)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("it should return ok when task update done successfully.", func(t *testing.T) {
		topic := request.UpdateTask{
			Title:       "Complete Integration Test",
			Description: "Kindly Complete Integration Testing of Task Module.",
		}
		jsonValue, _ := json.Marshal(topic)

		client := &http.Client{}
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/tasks/%v", ts.URL, id), bytes.NewBuffer(jsonValue))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		_, _ = buf.ReadFrom(resp.Body)
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		fmt.Println(responseBody)
		assert.Equal(t, 200, resp.StatusCode)
	})

	utils.ClearMockData(dbConn)
}

func TestTeamCRUD(t *testing.T) {
	buf := new(bytes.Buffer)
	var id int64 = 0
	ts, dbConn := runTestServer()
	defer ts.Close()

	t.Run("it should return ok when team creation done successfully.", func(t *testing.T) {
		team := request.CreateTeam{
			Details: request.Team{
				Name: "Team Jupiter",
			},
		}
		jsonValue, err := json.Marshal(team)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/teams", ts.URL), bytes.NewBuffer(jsonValue))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		var response response.SuccessResponse
		json.Unmarshal([]byte(responseBody), &response)
		fmt.Println(responseBody)
		id = *response.ID
		assert.Equal(t, constant.TEAM_CREATED, response.Message)
	})

	t.Run("it should return ok when team members added successfully.", func(t *testing.T) {
		teamMembers := request.TeamMembersWithTeamID{
			MemberIDs: []int64{
				954497896847212545,
			},
		}
		jsonValue, err := json.Marshal(teamMembers)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/teams/%v/members", ts.URL, id), bytes.NewBuffer(jsonValue))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		var response response.SuccessResponse
		json.Unmarshal([]byte(responseBody), &response)
		fmt.Println(responseBody)
		assert.Equal(t, constant.MEMBERS_ADDED_TO_TEAM, response.Message)
	})

	t.Run("it should return ok when team members removed successfully.", func(t *testing.T) {
		teamMembers := request.TeamMembersWithTeamID{
			TeamID: id,
			MemberIDs: []int64{
				954497896847212545,
			},
		}
		jsonValue, err := json.Marshal(teamMembers)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/teams/%v/members", ts.URL, id), bytes.NewBuffer(jsonValue))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		var response response.SuccessResponse
		json.Unmarshal([]byte(responseBody), &response)
		fmt.Println(responseBody)
		assert.Equal(t, constant.MEMBERS_REMOVED_FROM_TEAM, response.Message)
	})

	t.Run("it should return teams created by me.", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/teams?createdByMe=false", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		buf.ReadFrom(resp.Body)
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		fmt.Println(responseBody)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("it should return teams in which i were added.", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/teams?createdByMe=true", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		buf.Reset()
		buf.ReadFrom(resp.Body)
		responseBody := strings.TrimSuffix(buf.String(), "\n")
		fmt.Println(responseBody)
		assert.Equal(t, 200, resp.StatusCode)
	})
	
	utils.ClearMockData(dbConn)

}
