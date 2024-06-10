CREATE TABLE image_metadata (
    image_id TEXT PRIMARY KEY,
    upload_complete BOOLEAN NOT NULL DEFAULT FALSE,
    account_id INT NOT NULL REFERENCES users(account_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
;
