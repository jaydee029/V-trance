package main

import (
	"V-trance/user/internal/api"
	"V-trance/user/internal/database"
	"V-trance/user/internal/middleware"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load(".env")

	logger, _ := zap.NewProduction()

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		logger.Fatal("JWT secret key not set")
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

	h := api.New(jwt_secret, queries, logger)

	port := os.Getenv("PORT")

	r := chi.NewRouter()
	s := chi.NewRouter()
	r.Mount("user", s)

	s.Post("/signup", h.CreateUser)
	s.Post("/login", h.UserLogin)
	s.Post("/refresh", h.VerifyRefresh)
	s.Post("/revoke", h.RevokeToken)

	sermux := middleware.Corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The server is live on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
