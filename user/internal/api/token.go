package api

import (
	"V-trance/user/internal/auth"
	"V-trance/user/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (cfg *Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "bytes couldn't be converted")
		return
	}

	var pgtime pgtype.Timestamp

	err = pgtime.Scan(time.Now().UTC())
	if err != nil {
		cfg.logger.Info("Error setting timestamp:", zap.Error(err))
	}

	err = cfg.DB.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     []byte(token),
		RevokedAt: pgtime,
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	respondWithJson(w, http.StatusOK, "Token Revoked")
}

func (cfg *Handler) VerifyRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	is_refresh, err := auth.VerifyRefresh(token, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if !is_refresh {
		respondWithError(w, http.StatusUnauthorized, "Header doesn't contain refresh token")
		return
	}

	is_revoked, err := cfg.DB.VerifyRefresh(r.Context(), []byte(token))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if is_revoked {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token revoked")
		return
	}
	Idstr, err := auth.ValidateToken(token, cfg.Jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Id, err := uuid.Parse(Idstr)
	if err != nil {
		cfg.logger.Info("Error parsing string to UUID:", zap.Error(err))
	}

	auth_token, err := auth.Tokenize(Id, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, Token{
		Token: auth_token,
	})
}

func (cfg *Handler) VerifyUser(w http.ResponseWriter, r *http.Response) {
	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	useridstr, err := auth.ValidateToken(token, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var userid pgtype.UUID
	err = userid.Scan(useridstr)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, UserRes{
		ID: userid,
	})
}
