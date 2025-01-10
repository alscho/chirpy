-- name: GetHashedPasswordByEmail :one
SELECT * FROM users
WHERE email = $1;