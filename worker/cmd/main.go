package main

import (
	"V-trance/worker/Internal/database"
	"V-trance/worker/Internal/service"
	"context"
	"encoding/gob"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaydee029/V-trance/pubsub"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load("../.env")

	logger, _ := zap.NewProduction()

	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" {
		logger.Fatal("Object storage endpoint not set")
	}

	accessKeyID := os.Getenv("ACCESSKEYID")
	if accessKeyID == "" {
		logger.Fatal("AccessKeyID not set")
	}

	secretAccessKey := os.Getenv("SECRETACCESSKEY")
	if secretAccessKey == "" {
		logger.Fatal("SecretAccessKey not set")
	}

	bucketName := os.Getenv("BUCKET")
	if bucketName == "" {
		logger.Fatal("Bucket name not set")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		logger.Fatal("database connection string not set")
	}

	dbcon, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("Unable to connect to database:", zap.Error(err))
		os.Exit(1)
	}

	queries := database.New(dbcon)

	rabbitConnString := os.Getenv("RMQ_CONN")
	if rabbitConnString == "" {
		logger.Fatal("rabbitmq connection string not set")
	}

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	pb := pubsub.New(conn)
	gob.Register(&service.Task{})

	exchange := "vtrance-direct"
	key := "jobs"
	useSSL := true
	pathprefix := "users/"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	h := service.New(queries, pb, logger, minioClient, bucketName, pathprefix, exchange, key)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	log.Println("Worker service is live..")

	go h.EventListner(ctx, &wg)

	<-stop
	log.Println("Worker shutting down...")
	cancel() // cancel semaphore context
	wg.Wait()
	log.Println("All tasks completed.")

}
