-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS microInteract (
 id SERIAL PRIMARY KEY,
 name VARCHAR(255) NOT NULL,
 age VARCHAR(255),
 gender VARCHAR(255)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS microInteract;
-- +goose StatementEnd
