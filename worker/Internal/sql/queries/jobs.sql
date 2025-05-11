-- name: FetchJob :one
SELECT Video_id,Name,Type,Options From jobs WHERE Job_id=$1 AND Status=$2;

-- name: SetStatusJob :one
UPDATE jobs SET Status=$1 WHERE Job_id=$2
RETURNING *;