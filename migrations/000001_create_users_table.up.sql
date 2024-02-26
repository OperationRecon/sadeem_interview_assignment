CREATE TABLE IF NOT EXISTS users (
    id bigserail PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
);