CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
  username VARCHAR(50) UNIQUE,
  password TEXT,
  email VARCHAR(50) UNIQUE,
  created_on timestamptz NOT NULL,
  last_login timestamptz NULL
);