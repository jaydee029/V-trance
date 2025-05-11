-- name: FetchVideo :one
SELECT User_id,Name,Video_url,Resolution FROM videos WHERE Video_id=$1;

-- name: InsertVideoUrl :one
UPDATE videos SET Video_url=$1 WHERE Video_id=$2
RETURNING Video_id;