-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + interval '60' day
)
RETURNING *;

-- name: GetRefreshToken :one
select * from refresh_tokens where token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
where token = $1;