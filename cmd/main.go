package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chirag1807/task-management-system/api/route"
	"github.com/chirag1807/task-management-system/utils/socket"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/db"
)

func main() {
	config.LoadConfig("../.config", "../.config/secret.json")
	dbConn, redisClient, err := db.SetDBConection(0)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(context.Background())

	socketServer := socket.SocketConnection()

	r := route.InitializeRouter(dbConn, redisClient, socketServer)

	port := fmt.Sprintf(":%d", config.Config.Port)
	srv := &http.Server{
		Addr:        "localhost" + port,
		Handler:     r,
		IdleTimeout: 2 * time.Minute,
	}

	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
		log.Println("yes")
		defer socketServer.Close()
	}()

	log.Println("Server Started on Port", port)
	log.Fatal(srv.ListenAndServe())
}
