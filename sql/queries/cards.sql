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

-- name: CardsReset :exec
DELETE FROM cards;

-- name: GetCard :one
SELECT * FROM cards WHERE id = $1;

-- name: GetCardByName :one
SELECT * FROM cards WHERE NAME = $1;

-- name: GetCardsByRarity :many
SELECT * FROM cards WHERE rarity = $1;

-- name: GetCardsByType :many
SELECT * FROM cards WHERE type = $1;

-- name: GetCardNumber :one
SELECT COUNT(*) from cards;