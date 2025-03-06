-- Create user_characters table
CREATE TABLE user_characters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Character identifiers
    character_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    realm VARCHAR(100) NOT NULL,
    region VARCHAR(5) NOT NULL,
    
    -- Basic information
    class VARCHAR(50) NOT NULL,
    race VARCHAR(50) NOT NULL,
    gender VARCHAR(20) NOT NULL,
    faction VARCHAR(20) NOT NULL,
    active_spec_name VARCHAR(50) NOT NULL,
    active_spec_id INTEGER NOT NULL,
    active_spec_role VARCHAR(20) NOT NULL,
    
    -- Level and main stats
    level INTEGER NOT NULL DEFAULT 80,
    item_level DECIMAL(6,2) NOT NULL,
    mythic_plus_rating DECIMAL(6,2) NULL,
    mythic_plus_rating_color VARCHAR(10) NULL,
    achievement_points INTEGER NOT NULL DEFAULT 0,
    honorable_kills INTEGER NOT NULL DEFAULT 0,
    
    -- Image URLs
    avatar_url TEXT,
    inset_avatar_url TEXT,
    main_raw_url TEXT,
    profile_url TEXT,
    
    -- JSON structured data
    equipment_json JSONB,
    stats_json JSONB,
    talents_json JSONB,
    mythic_plus_json JSONB,
    raids_json JSONB,
    
    -- Metadata
    is_displayed BOOLEAN DEFAULT TRUE,
    last_api_update TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    
    -- Constraints and indexes
    UNIQUE(character_id, realm, region)
);

-- Create indexes to improve query performance
CREATE INDEX idx_user_characters_user_id ON user_characters(user_id);
CREATE INDEX idx_user_characters_character_id ON user_characters(character_id);
CREATE INDEX idx_user_characters_realm_region ON user_characters(realm, region);
CREATE INDEX idx_user_characters_deleted_at ON user_characters(deleted_at);
CREATE INDEX idx_user_characters_is_displayed ON user_characters(is_displayed) WHERE is_displayed = TRUE;
CREATE INDEX idx_user_characters_item_level ON user_characters(item_level);
CREATE INDEX idx_user_characters_mythic_plus_rating ON user_characters(mythic_plus_rating);

-- Add favorite_character_id column to users table
ALTER TABLE users ADD COLUMN favorite_character_id INTEGER NULL;
ALTER TABLE users ADD CONSTRAINT fk_users_favorite_character FOREIGN KEY (favorite_character_id) REFERENCES user_characters(id) ON DELETE SET NULL;