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
SET password = ?
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
