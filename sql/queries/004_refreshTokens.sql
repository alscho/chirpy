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