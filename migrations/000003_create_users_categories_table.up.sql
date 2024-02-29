CREATE TABLE IF NOT EXISTS user_categories (
    user_id bigserial NOT NULL REFERENCES users(id),
    category_id bigserial UNIQUE NOT NULL REFERENCES categories(id)
);