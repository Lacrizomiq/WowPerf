-- Migration DOWN: Remove reset password fields
-- Description: Remove columns added for reset password functionality
-- Version: 13
-- Created at: 2025-01-06

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS idx_users_reset_password_token,
    DROP COLUMN IF EXISTS reset_password_expires,
    DROP COLUMN IF EXISTS reset_password_token;