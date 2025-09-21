-- name: CreateCard :one
INSERT INTO cards (id, created_at, updated_at, name, rarity, type)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;