CREATE TABLE image_metadata (
    image_id TEXT PRIMARY KEY,
    image_name TEXT UNIQUE NOT NULL,
    upload_complete BOOLEAN NOT NULL DEFAULT FALSE,
    account_id INT NOT NULL REFERENCES users(account_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
;
