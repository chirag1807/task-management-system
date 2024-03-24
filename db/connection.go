package db

import (
	"context"
	"log"

	"github.com/chirag1807/task-management-system/config"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

func SetDBConection(flag int) (*pgx.Conn, *redis.Client, error) {
	//flag = 0 => main database connection, flag = 1 => test database connection
	var connConfig *pgx.ConnConfig
	var err error
	if flag == 0 {
		connConfig, err = pgx.ParseConfig(dbConnString())
	} else {
		connConfig, err = pgx.ParseConfig(testDbConnString())
	}
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	dbConn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Port,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return dbConn, redisClient, nil
}

func dbConnString() string {
	return "postgresql://root@127.0.0.1:26257/taskmanager?sslmode=disable"
}

func testDbConnString() string {
	return "postgresql://root@127.0.0.1:26257/testtaskmanager?sslmode=disable"
}