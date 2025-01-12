-- name: UpgradeUserToChirpyRed :execresult
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;