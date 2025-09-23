-- name: CreateRarity :one
INSERT INTO raritylist (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: RaritylistReset :exec
DELETE FROM raritylist;