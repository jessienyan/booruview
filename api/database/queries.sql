-- name: GetUser :one
SELECT * FROM users WHERE LOWER(username) = LOWER(@username);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :one
INSERT INTO users
(username, password, password_salt)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: DeleteUserData :exec
DELETE FROM user_data WHERE user_id = ?;

-- name: UserLoggedIn :exec
UPDATE users
SET last_login = ?
WHERE id = ?;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = ?, password_changed_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: GetUserData :one
SELECT * FROM user_data
WHERE user_id = ?;

-- name: CreateUserData :one
INSERT INTO user_data
(user_id, data)
VALUES (?, ?)
RETURNING *;

-- name: UpdateUserData :exec
UPDATE user_data
SET data = ?, updated_at = CURRENT_TIMESTAMP
WHERE user_id = ?;

-- name: CreateSession :one
INSERT INTO user_sessions (key, user_id, expires_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSessionByKey :one
SELECT * FROM user_sessions WHERE key = ?;

-- name: DeleteSessionByKey :exec
DELETE FROM user_sessions WHERE key = ?;

-- name: DeleteUserSessions :exec
DELETE FROM user_sessions WHERE user_id = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM user_sessions WHERE expires_at < CURRENT_TIMESTAMP;
