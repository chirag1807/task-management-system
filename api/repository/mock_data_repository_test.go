package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/db"
	"github.com/chirag1807/task-management-system/utils/socket"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn
var redisClient *redis.Client
var socketServer *socketio.Server

func init() {
	config.LoadConfig("../../.config/", "../../.config/secret.json")
	dbConn, redisClient, _ = db.SetDBConection(1)
	socketServer = socket.SocketConnection()
}

func TestMain(m *testing.M) {
	err := ClearMockData(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := dbConn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		log.Fatal(err)
	}

	tx, err = InsertMockData(tx)
	if err != nil {
		tx.Rollback(context.Background())
		log.Fatal(err)
	}
	tx.Commit(context.Background())

	// this is for running test of the controller
	//so from here it will go to actual function.
	exitVal := m.Run()

	err = ClearMockData(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func InsertMockData(tx pgx.Tx) (pgx.Tx, error) {
	batch := &pgx.Batch{}
	batch.Queue("INSERT INTO users (first_name, last_name, bio, email, password, profile) VALUES('Aashutosh', 'Gupta', 'Junior Software Engineer at ZURU TECH INDIA', 'guptaaahutosh354@gmail.com', '$2a$14$FhDiMSnCN8sJ7Tb0UDBXn.bbKVYF3b4ZVwEwPXfAzvDgXZlC3B1g2', 'Public');")
	batch.Queue("INSERT INTO tasks (title, description, deadline, assignee_team, status, priority, created_by, created_at) VALUES('task2', 'this is task2', '2024-03-30T22:59:59.000Z', 954507580144451585, 'TO-DO', 'Very High', 954488202459119617, current_timestamp());")
	results := tx.SendBatch(context.Background(), batch)
	defer results.Close()

	if err := results.Close(); err != nil {
		tx.Rollback(context.Background())
		log.Fatal(err)
		return tx, err
	}
	return tx, nil
}

func ClearMockData(dbConn *pgx.Conn) error {
	query := "DELETE FROM tasks WHERE id <> 954511608047501313;" +
		"DELETE FROM team_members WHERE NOT (team_id = 954507580144451585 AND member_id = 954488202459119617);" +
		"DELETE FROM teams WHERE id <> 954507580144451585;" +
		"DELETE FROM refresh_tokens WHERE user_id <> 954488202459119617;" +
		"DELETE FROM users WHERE id NOT IN (954488202459119617, 954497896847212545);" +
		"DELETE FROM otps WHERE id <> 954537852771565569;"

	_, err := dbConn.Exec(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
