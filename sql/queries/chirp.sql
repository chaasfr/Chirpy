-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: GetAllChirp :many
SELECT * from chirps order by created_at asc;

-- name: GetChirp :one
SELECT * from chirps where id = $1;