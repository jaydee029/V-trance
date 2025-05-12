-- name: CreateJob :one
INSERT INTO jobs(Job_id,Video_id,Name,Type, Options, Status, Created_at) VALUES($1,$2,$3,$4,$5,$6,$7)
RETURNING Name, Job_id, Video_id;

-- name: FetchStatus :one
SELECT Name, Video_id, Status FROM jobs WHERE Job_id=$1;

-- name: FetchJob :one
SELECT Video_id,Name,Type,Options From jobs WHERE Job_id=$1 AND Status=$2;

-- name: SetStatusJob :one
UPDATE jobs SET Status=$1 WHERE Job_id=$2
RETURNING *;