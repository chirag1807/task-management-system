package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chirag1807/task-management-system/api/route"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/db"
	"github.com/chirag1807/task-management-system/docs"
	"github.com/chirag1807/task-management-system/utils/socket"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Task Manager API Documentation
// @version 1.0
// @description This is the api documentation of task manager project.
// @host localhost:9090
// @BasePath /
// @query.collection.format multi
func main() {
	config.LoadConfig("../.config", "../.config/secret.json")
	dbConn, redisClient, rabbitmqConn := db.SetDBConection(0)
	defer dbConn.Close(context.Background())
  
	socketServer := socket.SocketConnection()
	port := fmt.Sprintf("localhost:%d", config.Config.Port)

	docs.SwaggerInfo.Title = "Task Manager API Documentation"
	docs.SwaggerInfo.Description = "This is a swagger demo api documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost" + port
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := route.InitializeRouter(dbConn, redisClient, rabbitmqConn, socketServer)
	r.Mount("/swagger", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:        port,
		Handler:     r,
		IdleTimeout: 2 * time.Minute,
	}

	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
		defer socketServer.Close()
	}()

	log.Println("Server Started on Port " + port)
	log.Fatal(srv.ListenAndServe())
}
