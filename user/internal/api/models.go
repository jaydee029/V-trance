package api

import "github.com/jackc/pgx/v5/pgtype"

type Token struct {
	Token string `json:"token"`
}
type LoginInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
type UserRes struct {
	ID       pgtype.UUID `json:"id"`
	Email    string      `json:"email,omitempty"`
	Name     string      `json:"name"`
	Username string      `json:"username,omitempty"`
}
type LoginRes struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
}
type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
