CREATE TABLE IF NOT EXISTS categories (
    id bigserial PRIMARY key,
    name text UNIQUE not NULL
);

-- sample data
INSERT INTO categories (name) values ('cat1');
INSERT INTO categories (name) values ('thing');
INSERT INTO categories (name) values ('something');