-- Drop the new index
DROP INDEX IF EXISTS users_battle_net_id_key;

-- Restore battle_net_id to integer
ALTER TABLE users 
    ALTER COLUMN battle_net_id TYPE INTEGER USING 
        CASE 
            WHEN battle_net_id ~ '^[0-9]+$' THEN battle_net_id::INTEGER
            ELSE NULL
        END;

-- Add back unique constraint
ALTER TABLE users 
    ADD CONSTRAINT users_battle_net_id_key UNIQUE (battle_net_id);

-- Rename encrypted_access_token back to encrypted_token
ALTER TABLE users 
    RENAME COLUMN encrypted_access_token TO encrypted_token;

-- Add back battle_net_refresh_token
ALTER TABLE users 
    ADD COLUMN battle_net_refresh_token TEXT;

-- Remove new columns
ALTER TABLE users 
    DROP COLUMN IF EXISTS encrypted_refresh_token,
    DROP COLUMN IF EXISTS last_token_refresh;