-- Migration UP: Add reset password fields
-- Description: Add columns for reset password token and expiration
-- Version: 13
-- Created at: 2025-01-06

DO $$
BEGIN
    -- Add reset_password_token column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'reset_password_token'
    ) THEN
        ALTER TABLE users ADD COLUMN reset_password_token VARCHAR(255) NULL;
    END IF;
    
    -- Add reset_password_expires column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'reset_password_expires'
    ) THEN
        ALTER TABLE users ADD COLUMN reset_password_expires TIMESTAMP WITH TIME ZONE NULL;
    END IF;
    
    -- Add constraint if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'idx_users_reset_password_token'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT idx_users_reset_password_token UNIQUE (reset_password_token);
    END IF;
END
$$;