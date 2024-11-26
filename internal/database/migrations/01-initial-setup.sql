-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
  id TEXT PRIMARY KEY,
  username TEXT UNIQUE,
  password_hash TEXT NOT NULL,
  roles TEXT[] NOT NULL
);

CREATE TABLE user_accounts (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  service TEXT NOT NULL,
  username TEXT NOT NULL,
  token TEXT NOT NULL,

  refresh_started_at TIMESTAMP DEFAULT NULL,
  refresh_active_at TIMESTAMP DEFAULT NULL,
  refresh_completed_at TIMESTAMP DEFAULT NULL,
  refresh_status TIMESTAMP DEFAULT NULL
);

CREATE TABLE decks (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  service TEXT NOT NULL,
  service_id TEXT NOT NULL,
  name TEXT NOT NULL,
  format TEXT NOT NULL,
  url TEXT NOT NULL,
  color_identity TEXT[] NOT NULL,
  folder JSONB NOT NULL DEFAULT '{}',
  leaders JSONB NOT NULL DEFAULT '{}',
  archetypes JSONB NOT NULL DEFAULT '[]',

  updated_at TIMESTAMP DEFAULT NULL,
  refreshed_at TIMESTAMP DEFAULT NULL
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE decks;
DROP TABLE user_accounts;
DROP TABLE users;
