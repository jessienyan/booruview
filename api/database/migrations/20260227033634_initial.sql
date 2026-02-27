-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login      DATETIME,
    username        VARCHAR(16) NOT NULL,
    password        BLOB NOT NULL,
    password_salt   BLOB NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS uniq_users_username ON users(LOWER(username));
CREATE TABLE IF NOT EXISTS user_data (
    user_id     INTEGER PRIMARY KEY,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    data        TEXT NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_data;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
