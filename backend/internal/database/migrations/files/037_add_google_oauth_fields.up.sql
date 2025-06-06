-- Add Google OAuth fields to users table
-- Migration 037: Google OAuth Integration

-- Add Google OAuth columns
ALTER TABLE users ADD COLUMN google_id VARCHAR(255);
ALTER TABLE users ADD COLUMN google_email VARCHAR(255);

-- Add unique constraint on google_id (only if not null)
CREATE UNIQUE INDEX users_google_id_unique ON users(google_id) WHERE google_id IS NOT NULL;

-- Add performance indexes
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
CREATE INDEX idx_users_google_email ON users(google_email) WHERE google_email IS NOT NULL;