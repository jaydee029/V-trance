package service

import (
	"V-trance/worker/Internal/database"

	"github.com/jaydee029/V-trance/pubsub"
	"go.uber.org/zap"
)

type Handler struct {
	DB       *database.Queries
	Pb       *pubsub.PubSub
	logger   *zap.Logger
	Exchange string
	key      string
}

func New(db *database.Queries, pb *pubsub.PubSub, lg *zap.Logger, exchange, key string) *Handler {
	return &Handler{
		DB:       db,
		Pb:       pb,
		logger:   lg,
		Exchange: exchange,
		key:      key,
	}
}
