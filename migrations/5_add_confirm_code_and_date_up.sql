ALTER TABLE users ADD COLUMN confirm_code VARCHAR(255);
ALTER TABLE users ADD COLUMN verified_at timestamp with time zone;
