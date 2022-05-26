-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_notes (
                         id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
                         user_id uuid REFERENCES users(id) ON DELETE CASCADE,
                         type INTEGER,
                         title VARCHAR(50) NOT NULL,
                         body bytea NOT NULL,
                         created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_notes;
-- +goose StatementEnd