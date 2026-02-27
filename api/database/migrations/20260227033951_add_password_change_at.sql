-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN password_changed_at DATETIME;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN password_changed_at;
-- +goose StatementEnd
