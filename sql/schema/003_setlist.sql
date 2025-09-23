-- +goose Up
CREATE TABLE setlist (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name text UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE setlist;