package db

import (
	"context"
	"log"

	"github.com/chirag1807/task-management-system/config"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetDBConection(flag int) (*pgx.Conn, *redis.Client, *amqp.Connection) {
	//flag = 0 => main database connection, flag = 1 => test database connection
	var connConfig *pgx.ConnConfig
	var err error
	if flag == 0 {
		connConfig, err = pgx.ParseConfig(dbConnString())
	} else {
		connConfig, err = pgx.ParseConfig(testDbConnString())
	}
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Port,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	rabbitmqConn, err := amqp.Dial("amqp://" + config.Config.RabbitMQ.Username + ":" + config.Config.RabbitMQ.Password + "@localhost:" + config.Config.RabbitMQ.Port + "/")
	if err != nil {
		log.Fatal(err)
	}

	return dbConn, redisClient, rabbitmqConn
}

func dbConnString() string {
	return "postgresql://" + config.Config.Database.Username + ":" + config.Config.Database.Password + "@cockroachdb:" + config.Config.Database.Port + "/" + config.Config.Database.Name + "?sslmode=" + config.Config.Database.SSLMode
}

func testDbConnString() string {
	return "postgresql://" + config.Config.Database.Username + ":" + config.Config.Database.Password + "@cockroachdb:" + config.Config.Database.Port + "/" + config.Config.Database.TestDatabaseName + "?sslmode=" + config.Config.Database.SSLMode
}
