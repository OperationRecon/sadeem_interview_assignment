CREATE TABLE IF NOT EXISTS admins (
    id bigserial NOT NULL REFERENCES users(id)
);