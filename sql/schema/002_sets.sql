-- +goose Up
CREATE TABLE sets (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name text NOT NULL,
    card_id uuid NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    UNIQUE (name, card_id)
);

-- +goose Down
DROP TABLE sets;