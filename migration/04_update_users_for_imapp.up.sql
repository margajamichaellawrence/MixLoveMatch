-- Add email column
ALTER TABLE users ADD COLUMN email VARCHAR(255) UNIQUE;

-- Note: ID type change from BIGINT to VARCHAR(36) for UUID requires recreation
-- For now, we'll keep BIGINT but in production you'd want to:
-- 1. Create new table with VARCHAR ID
-- 2. Migrate data
-- 3. Drop old table
-- 4. Rename new table
