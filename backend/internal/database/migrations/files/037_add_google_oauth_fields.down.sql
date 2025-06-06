-- Rollback Google OAuth fields from users table
-- Migration 037: Google OAuth Integration - ROLLBACK

-- Drop indexes first
DROP INDEX IF EXISTS idx_users_google_email;
DROP INDEX IF EXISTS idx_users_google_id;
DROP INDEX IF EXISTS users_google_id_unique;

-- Drop columns
ALTER TABLE users DROP COLUMN IF EXISTS google_email;
ALTER TABLE users DROP COLUMN IF EXISTS google_id;