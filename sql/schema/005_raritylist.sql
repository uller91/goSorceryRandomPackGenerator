-- +goose Up
CREATE TABLE raritylist (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name text UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE raritylist;