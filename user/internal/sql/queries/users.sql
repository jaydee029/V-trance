-- name: CreateUser :one
INSERT INTO users(name,Email,passwd,id,created_at,username) VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetUserEmail :one
SELECT * FROM users WHERE Email=$1;

-- name: GetUserUsername :one
SELECT * FROM users WHERE username=$1;

-- name: Is_Email :one
SELECT EXISTS (SELECT 1 FROM users WHERE Email=$1) AS value_exists;

-- name: Is_Username :one
SELECT EXISTS (SELECT 1 FROM users WHERE username=$1) AS value_exists;