-- WARNING: this file is ran every time the server starts
-- Follow these rules when making changes:
--  1. Add your changes to the end of the file
--  2. Do NOT edit the existing SQL
--  3. Wrap your changes in a transaction
--  4. Your changes MUST be idempotent


-- Initial schema
BEGIN;
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
COMMIT;


-- 2026-02-26
-- Track when users last changed their password in order to expire existing auth tokens.
-- If null, the user has never changed their password since registering their account.
BEGIN;
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_changed_at DATETIME;
COMMIT;
