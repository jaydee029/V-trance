package api

import (
	"V-trance/user/internal/database"

	"go.uber.org/zap"
)

type Handler struct {
	Jwtsecret string
	DB        *database.Queries
	logger    *zap.Logger
}

func New(jwt string, DBQueries *database.Queries, logger *zap.Logger) *Handler {
	return &Handler{
		Jwtsecret: jwt,
		DB:        DBQueries,
		logger:    logger,
	}
}
