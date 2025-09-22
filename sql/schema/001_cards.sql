-- +goose Up
CREATE TABLE cards (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name text UNIQUE NOT NULL,
    rarity text NOT NULL,
    type text NOT NULL
);

-- +goose Down
DROP TABLE cards;