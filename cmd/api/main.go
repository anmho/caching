package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/anmho/caching/api"
	"github.com/anmho/caching/cache"
	"github.com/anmho/caching/todo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
)

const (
	port = 8080

	RedisURL      = "http://localhost:6380"
	RedisPassword = "password"

	DynamoDBURL = "http://localhost:8000"
)

func WithEndpoint(endpoint string) func(*dynamodb.Options) {
	return func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     RedisURL,
		Password: RedisPassword,
	})
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	// Setup dependencies
	dynamoClient := dynamodb.NewFromConfig(cfg, WithEndpoint(DynamoDBURL))
	todoCache := cache.New[todo.Todo](redisClient)
	todoService := todo.MakeService(
		dynamoClient,
		todoCache,
	)

	mux := api.New(todoService)
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("listening on port %d\n", port)
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// shutdown
			log.Println("shutting down mux")
			os.Exit(0)
		}
	}
}
