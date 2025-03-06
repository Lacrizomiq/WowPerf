-- Remove constraint and column from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_favorite_character;
ALTER TABLE users DROP COLUMN IF EXISTS favorite_character_id;

-- Remove indexes
DROP INDEX IF EXISTS idx_user_characters_user_id;
DROP INDEX IF EXISTS idx_user_characters_character_id;
DROP INDEX IF EXISTS idx_user_characters_realm_region;
DROP INDEX IF EXISTS idx_user_characters_deleted_at;
DROP INDEX IF EXISTS idx_user_characters_is_displayed;
DROP INDEX IF EXISTS idx_user_characters_item_level;
DROP INDEX IF EXISTS idx_user_characters_mythic_plus_rating;

-- Drop user_characters table
DROP TABLE IF EXISTS user_characters;