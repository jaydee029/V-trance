-- name: FetchVideo :one
SELECT Name,Video_url FROM videos WHERE Video_id=$1;