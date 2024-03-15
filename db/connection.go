package db

import (
	"context"
	"log"

	"github.com/chirag1807/task-management-system/config"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

func SetDBConection() (*pgx.Conn, *redis.Client, error) {
	connConfig, err := pgx.ParseConfig(dbConnString())
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
	return "postgresql://" + config.Config.Database.Username + ":" + config.Config.Database.Password + "@127.0.0.1:" + config.Config.Database.Port + "/" + config.Config.Database.Name + "?sslmode=" + config.Config.Database.SSLMode
}
