-- name: InsertInitialDetails :one
INSERT INTO videos(User_id,Video_id,Name,Type,Resolution, Created_at) VALUES ($1,$2,$3,$4,$5,$6)
RETURNING Name,Video_id;

-- name: IfVideoExists :one
SELECT EXISTS (SELECT 1 FROM videos WHERE Video_id=$1) AS value_exists;

-- name: InsertFinalVideoDetails :one
UPDATE videos SET Video_url=$1 WHERE Video_id=$2 RETURNING Name, Video_id;

-- name: GetVideos :many
SELECT Name, Stream_url FROM videos WHERE User_id=$1;

-- name: GetStreamurl :one
SELECT Name, Stream_url FROM videos WHERE User_id=$1 AND Video_id=$2;