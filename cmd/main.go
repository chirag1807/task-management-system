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
)

func main() {
	config.LoadConfig("../.config")
	dbConn, redisClient, err := db.SetDBConection()
	if err != nil {
		log.Fatal(err)
		//handle error here
	}
	defer dbConn.Close(context.Background())

	port := fmt.Sprintf(":%d", config.Config.Port)
	srv := &http.Server{
		Addr:        port,
		Handler:     route.InitializeRouter(dbConn, redisClient),
		IdleTimeout: 2 * time.Minute,
	}

	log.Println("Server Started on Port", port)
	log.Fatal(srv.ListenAndServe())
}
