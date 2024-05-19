package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/db"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/chirag1807/task-management-system/utils/socket"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

var dbConn *pgx.Conn
var redisClient *redis.Client
var rabbitmqConn *amqp.Connection
var socketServer *socketio.Server

func init() {
	config.LoadConfig("../../.config/", "../../.config/secret.json")
	dbConn, redisClient, rabbitmqConn = db.SetDBConection(1)
	socketServer = socket.SocketConnection()
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
