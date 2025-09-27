-- name: CreateSetAndCard :one
INSERT INTO sets (id, created_at, updated_at, name, card_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: SetsReset :exec
DELETE FROM sets;

-- name: GetCardsBySet :many
SELECT * FROM sets WHERE name = $1;
