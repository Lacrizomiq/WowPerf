-- Migration UP: Add reset password fields
-- Description: Add columns for reset password token and expiration
-- Version: 13
-- Created at: 2025-01-06

ALTER TABLE users 
    ADD COLUMN reset_password_token VARCHAR(255) NULL,
    ADD COLUMN reset_password_expires TIMESTAMP WITH TIME ZONE NULL,
    ADD CONSTRAINT idx_users_reset_password_token UNIQUE (reset_password_token);
