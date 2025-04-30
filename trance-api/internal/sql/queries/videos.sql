-- name: GetStreamURL :one
SELECT stream_url FROM videos WHERE Video_id=$1;