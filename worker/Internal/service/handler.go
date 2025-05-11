package service

import (
	"V-trance/worker/Internal/database"

	"github.com/jaydee029/V-trance/pubsub"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type Handler struct {
	DB         *database.Queries
	Pb         *pubsub.PubSub
	logger     *zap.Logger
	B2Client   *minio.Client
	Bucket     string
	Pathprefix string
	Exchange   string
	key        string
}

func New(db *database.Queries, pb *pubsub.PubSub, lg *zap.Logger, b2client *minio.Client, bucket, pathprefix, exchange, key string) *Handler {
	return &Handler{
		DB:         db,
		Pb:         pb,
		logger:     lg,
		Exchange:   exchange,
		key:        key,
		B2Client:   b2client,
		Bucket:     bucket,
		Pathprefix: pathprefix,
	}
}
