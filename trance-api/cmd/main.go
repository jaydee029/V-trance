package main

import (
	"V-trance/trance-api/internal/api"
	"V-trance/trance-api/internal/database"
	"V-trance/trance-api/internal/middleware"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load(".env")

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
	useSSL := true
	pathprefix := "users/"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	handler := api.New(minioClient, bucketName, pathprefix, queries, dbcon, logger)

	port := os.Getenv("PORT")

	r := chi.NewRouter()
	s := chi.NewRouter()
	r.Mount("/tranceapi", s)

	s.Post("/upload-Url", handler.GetUploadUrl)
	s.Post("/notifyUpload", handler.NotifyUpload)
	s.Get("/getVideos", handler.GetVideos)
	s.Get("/jobStatus/{jobid}", handler.GetStatus)
	s.Get("/fetchVideo/{videoid}", handler.GetStreamUrl)
	s.Get("/downloadVideo/{videoid}", handler.GetDownloadUrl)

	sermux := middleware.Corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The server is live on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
