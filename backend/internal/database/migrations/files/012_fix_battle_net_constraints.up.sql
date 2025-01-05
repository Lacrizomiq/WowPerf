-- 012_fix_battle_net_constraints.up.sql

-- 1. Supprimer les anciennes contraintes
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_battle_tag_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_battle_net_id_key;

-- 2. Allow NULL for Battle.net fields
ALTER TABLE users 
    ALTER COLUMN battle_tag DROP NOT NULL,
    ALTER COLUMN battle_net_id DROP NOT NULL;

-- 3. Set empty values to NULL
UPDATE users 
SET battle_tag = NULL 
WHERE battle_tag = '';

UPDATE users 
SET battle_net_id = NULL 
WHERE battle_net_id = '';

-- 4. Create new constraints that ignore NULL
CREATE UNIQUE INDEX users_battle_tag_unique 
    ON users(battle_tag) 
    WHERE battle_tag IS NOT NULL 
    AND deleted_at IS NULL;

CREATE UNIQUE INDEX users_battle_net_id_unique 
    ON users(battle_net_id) 
    WHERE battle_net_id IS NOT NULL 
    AND deleted_at IS NULL;

-- 5. Add CHECK constraints to avoid empty strings
ALTER TABLE users 
    ADD CONSTRAINT chk_battle_tag_not_empty 
    CHECK (battle_tag IS NULL OR battle_tag <> '');

ALTER TABLE users 
    ADD CONSTRAINT chk_battle_net_id_not_empty 
    CHECK (battle_net_id IS NULL OR battle_net_id <> '');