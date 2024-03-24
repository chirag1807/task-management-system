package controller

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/chirag1807/task-management-system/api/repository"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/db"
	"github.com/chirag1807/task-management-system/utils/socket"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn
var rdb *redis.Client
var r *chi.Mux
var authService service.AuthService
var taskService service.TaskService
var teamService service.TeamService
var userService service.UserService

type contextKey string

var (
	TokenKey        = contextKey("token")
	UserIdKey       = contextKey("userId")
	SocketServerKey = contextKey("socketServer")
)

func init() {
	config.LoadConfig("../../.config/")
	dbConn, redisClient, _ := db.SetDBConection(1)
	socketServer := socket.SocketConnection()
	r = chi.NewRouter()

	authRepository := repository.NewAuthRepo(conn, rdb)
	authService = service.NewAuthService(authRepository)

	taskRepository := repository.NewTaskRepo(dbConn, redisClient, socketServer)
	taskService = service.NewTaskService(taskRepository)

	teamRepository := repository.NewTeamRepo(dbConn, redisClient)
	teamService = service.NewTeamService(teamRepository)

	userRepository := repository.NewUserRepo(dbConn, redisClient)
	userService = service.NewUserService(userRepository)
}

func TestMain(m *testing.M) {
	err := ClearMockData(conn)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}

	tx, err = InsertMockData(tx)
	if err != nil {
		tx.Rollback(context.Background())
		return
	}
	tx.Commit(context.Background())

	// this is for running test of the controller
	//so from here it will go to actual function.
	exitVal := m.Run()

	err = ClearMockData(conn)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func InsertMockData(tx pgx.Tx) (pgx.Tx, error) {
	_, err := tx.Exec(context.Background(), "INSERT INTO users (first_name, last_name, bio, email, password, profile) VALUES('Aashutosh', 'Gupta', 'Junior Software Engineer at ZURU TECH INDIA', 'guptaaahutosh354@gmail.com', 'Aashutosh123$', 'public');")
	if err != nil {
		return tx, err
	}
	return tx, nil
}

func ClearMockData(dbConn *pgx.Conn) error {
	query := "DELETE FROM tasks;" + "DELETE FROM team_members;" + "DELETE FROM teams;" + "DELETE FROM taskstatus;" +
		"DELETE FROM refresh_tokens;" + "DELETE FROM users;" + "DELETE FROM otps;"

	_, err := dbConn.Exec(context.Background(), query)
	if err != nil {
		log.Print(err)
	}
	return nil
}
