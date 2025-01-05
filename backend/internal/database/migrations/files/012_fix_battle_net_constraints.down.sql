-- 012_fix_battle_net_constraints.down.sql

-- 1. Delete CHECK constraints
ALTER TABLE users 
    DROP CONSTRAINT IF EXISTS chk_battle_tag_not_empty,
    DROP CONSTRAINT IF EXISTS chk_battle_net_id_not_empty;

-- 2. Delete unique indexes
DROP INDEX IF EXISTS users_battle_tag_unique;
DROP INDEX IF EXISTS users_battle_net_id_unique;

-- 3. Restore old constraints
-- Note: We first convert NULL to empty strings to avoid errors
UPDATE users 
SET battle_tag = '' 
WHERE battle_tag IS NULL;

UPDATE users 
SET battle_net_id = '' 
WHERE battle_net_id IS NULL;

ALTER TABLE users 
    ALTER COLUMN battle_tag SET NOT NULL,
    ALTER COLUMN battle_net_id SET NOT NULL;

-- 4. Recreate old unique constraints
ALTER TABLE users 
    ADD CONSTRAINT users_battle_tag_key UNIQUE (battle_tag);

ALTER TABLE users 
    ADD CONSTRAINT users_battle_net_id_key UNIQUE (battle_net_id);