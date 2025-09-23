-- name: CreateType :one
INSERT INTO typelist (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: TypelistReset :exec
DELETE FROM typelist;