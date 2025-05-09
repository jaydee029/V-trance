package api

import (
	"V-trance/trance-api/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type Handler struct {
	B2Client   *minio.Client
	Bucket     string
	Pathprefix string
	DB         *database.Queries
	DBPool     *pgxpool.Pool
	logger     *zap.Logger
}

func New(client *minio.Client, bucket, pathprefix string, db *database.Queries, dbpool *pgxpool.Pool, loggerclient *zap.Logger) *Handler {
	return &Handler{
		B2Client:   client,
		Bucket:     bucket,
		Pathprefix: pathprefix,
		DB:         db,
		DBPool:     dbpool,
		logger:     loggerclient,
	}
}
