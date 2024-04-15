package controller

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/chirag1807/logease"
	"github.com/chirag1807/task-management-system/api/repository"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/db"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/chirag1807/task-management-system/utils/socket"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

var dbConn *pgx.Conn
var redisClient *redis.Client
var rabbitmqConn *amqp.Connection
var r *chi.Mux
var socketServer *socketio.Server
var slogLoggerInstance logease.SlogLoggerInstance
var authService service.AuthService
var taskService service.TaskService
var teamService service.TeamService
var userService service.UserService

func init() {
	config.LoadConfig("../../.config/", "../../.config/secret.json")
	dbConn, redisClient, rabbitmqConn = db.SetDBConection(1)
	socketServer = socket.SocketConnection()
	loggerInstance, err := logease.InitLogease(false, config.Config.TeamsWebHookURL, logease.Slog)
	if err != nil {
		log.Fatal(err)
	}
	slogLoggerInstance = loggerInstance.(logease.SlogLoggerInstance)
	r = chi.NewRouter()

	authRepository := repository.NewAuthRepo(dbConn)
	authService = service.NewAuthService(authRepository)

	taskRepository := repository.NewTaskRepo(dbConn, redisClient, socketServer)
	taskService = service.NewTaskService(taskRepository)

	teamRepository := repository.NewTeamRepo(dbConn, redisClient)
	teamService = service.NewTeamService(teamRepository)

	userRepository := repository.NewUserRepo(dbConn, rabbitmqConn)
	userService = service.NewUserService(userRepository)
}

func TestMain(m *testing.M) {
	err := utils.ClearMockData(dbConn)
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

	// this is for running test of the controller
	//so from here it will go to actual function.
	exitVal := m.Run()

	err = utils.ClearMockData(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}