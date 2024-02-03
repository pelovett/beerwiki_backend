CREATE TABLE users (
    account_id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    user_name TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
)
;
