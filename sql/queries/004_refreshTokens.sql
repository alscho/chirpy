-- name: AddRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
);

-- name: GetUserFromValidRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE token = $1
    AND expires_at > NOW()
    AND revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: UpdatePasswordEmailFromToken :exec
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1;