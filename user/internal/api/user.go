package api

import (
	"V-trance/user/internal/auth"
	"V-trance/user/internal/database"
	"V-trance/user/internal/validate"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := UserInput{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	err = validate.ValidateEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	email_if_exist, err := cfg.DB.Is_Email(context.Background(), params.Email)

	if email_if_exist {
		respondWithError(w, http.StatusConflict, "Email already exists")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.ValidateUsername(params.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	username_if_exists, err := cfg.DB.Is_Username(r.Context(), params.Username)
	if username_if_exists {
		respondWithError(w, http.StatusConflict, "Username already exists")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.ValidatePassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	encrypted, _ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	uuids := uuid.New().String()
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(uuids)
	if err != nil {
		cfg.logger.Info("Error setting UUID:", zap.Error(err))
		return
	}

	var pgtime pgtype.Timestamp
	err = pgtime.Scan(time.Now().UTC())
	if err != nil {
		cfg.logger.Info("Error setting timestamp:", zap.Error(err))
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:      params.Name,
		Email:     params.Email,
		Passwd:    encrypted,
		ID:        pgUUID,
		CreatedAt: pgtime,
		Username:  params.Username,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, UserRes{
		Email:    user.Email,
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	})
}

func (cfg *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginInput{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	loginusername := true
	if params.Email != "" {
		loginusername = false
	}

	user := database.User{}

	if !loginusername {
		user, err = cfg.DB.GetUserEmail(r.Context(), params.Email)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

	} else {
		user, err = cfg.DB.GetUserEmail(r.Context(), params.Email)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password doesn't match")
		return
	}

	Userid, _ := uuid.FromBytes(user.ID.Bytes[:])

	Token, err := auth.Tokenize(Userid, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Refresh_token, err := auth.RefreshToken(Userid, cfg.Jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, LoginRes{
		Username:      user.Username,
		Email:         user.Email,
		Token:         Token,
		Refresh_token: Refresh_token,
	})

}
