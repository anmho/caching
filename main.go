package main

import (
	"context"
	"errors"
	"github.com/anmho/caching/api"
	"github.com/anmho/caching/todo"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"net/http"
	"os"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	// Setup dependencies
	dynamoClient := dynamodb.NewFromConfig(cfg)
	todoService := todo.MakeService(dynamoClient)

	todoAPI := api.New(todoService)
	srv := http.Server{
		Addr:    ":8080",
		Handler: todoAPI,
	}

	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// shutdown
			log.Println("shutting down server")
			os.Exit(0)
		}
	}
}
