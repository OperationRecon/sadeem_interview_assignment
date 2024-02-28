CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    pfp_filepath text DEFAULT 'defaultpfp.jpeg');