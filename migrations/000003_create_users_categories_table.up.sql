CREATE TABLE IF NOT EXISTS user_categories (
    user_id bigserial NOT NULL REFERENCES users(id),
    category_id bigserial NOT NULL REFERENCES categories(id)
);