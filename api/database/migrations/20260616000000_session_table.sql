-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_sessions (
	key			VARCHAR(32) PRIMARY KEY,
	created_at	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at	DATETIME NOT NULL,
	user_id		INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_sessions;
-- +goose StatementEnd
