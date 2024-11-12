-- 002_create_talent_system.up.sql

CREATE TABLE IF NOT EXISTS talent_trees (
    id SERIAL PRIMARY KEY,
    trait_tree_id INTEGER,
    spec_id INTEGER,
    class_name VARCHAR(255),
    class_id INTEGER,
    class_icon VARCHAR(255),
    spec_name VARCHAR(255),
    spec_icon VARCHAR(255),
    full_node_order INTEGER[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uni_talent_trees_trait_tree_spec_id UNIQUE (trait_tree_id, spec_id)
);

CREATE TABLE IF NOT EXISTS talent_nodes (
    id SERIAL PRIMARY KEY,
    node_id INTEGER,
    talent_tree_id INTEGER,
    spec_id INTEGER,
    node_type VARCHAR(50),
    name VARCHAR(255),
    type VARCHAR(50),
    pos_x FLOAT,
    pos_y FLOAT,
    max_ranks INTEGER,
    entry_node BOOLEAN,
    req_points INTEGER,
    free_node BOOLEAN,
    next INTEGER[],
    prev INTEGER[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uni_talent_nodes_node_tree_spec UNIQUE (node_id, talent_tree_id, spec_id)
);

CREATE TABLE IF NOT EXISTS talent_entries (
    id SERIAL PRIMARY KEY,
    entry_id INTEGER,
    node_id INTEGER,
    talent_tree_id INTEGER,
    spec_id INTEGER,
    definition_id INTEGER,
    max_ranks INTEGER,
    type VARCHAR(50),
    name VARCHAR(255),
    spell_id INTEGER,
    icon VARCHAR(255),
    index INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uni_talent_entries_entry_node_tree_spec UNIQUE (entry_id, node_id, talent_tree_id, spec_id)
);

CREATE TABLE IF NOT EXISTS sub_tree_nodes (
    id SERIAL PRIMARY KEY,
    sub_tree_node_id INTEGER UNIQUE,
    talent_tree_id INTEGER,
    spec_id INTEGER,
    name VARCHAR(255),
    type VARCHAR(50),
    pos_x FLOAT,
    pos_y FLOAT,
    entry_node BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uni_sub_tree_nodes_id_tree_spec UNIQUE (sub_tree_node_id, talent_tree_id, spec_id)
);

CREATE TABLE IF NOT EXISTS sub_tree_entries (
    id SERIAL PRIMARY KEY,
    sub_tree_node_id INTEGER,
    entry_id INTEGER,
    type VARCHAR(50),
    name VARCHAR(255),
    trait_sub_tree_id INTEGER,
    trait_tree_id INTEGER,
    atlas_member_name VARCHAR(255),
    nodes INTEGER[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS hero_nodes (
    id SERIAL PRIMARY KEY,
    node_id INTEGER,
    talent_tree_id INTEGER,
    spec_id INTEGER,
    name VARCHAR(255),
    type VARCHAR(50),
    pos_x FLOAT,
    pos_y FLOAT,
    max_ranks INTEGER,
    entry_node BOOLEAN,
    sub_tree_id INTEGER,
    requires_node INTEGER,
    next INTEGER[],
    prev INTEGER[],
    free_node BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uni_hero_nodes_node_tree_spec UNIQUE (node_id, talent_tree_id, spec_id)
);

CREATE TABLE IF NOT EXISTS hero_entries (
    id SERIAL PRIMARY KEY,
    entry_id INTEGER,
    node_id INTEGER,
    talent_tree_id INTEGER,
    spec_id INTEGER,
    definition_id INTEGER,
    max_ranks INTEGER,
    type VARCHAR(50),
    name VARCHAR(255),
    spell_id INTEGER,
    icon VARCHAR(255),
    index INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uni_hero_entries_entry_node_tree_spec UNIQUE (entry_id, node_id, talent_tree_id, spec_id)
);

-- Adding join table for sub_tree_node_talents
CREATE TABLE IF NOT EXISTS sub_tree_node_talents (
    sub_tree_node_id INTEGER REFERENCES sub_tree_nodes(id),
    talent_node_id INTEGER REFERENCES talent_nodes(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (sub_tree_node_id, talent_node_id)
);

-- Create indexes for soft delete
CREATE INDEX IF NOT EXISTS idx_talent_trees_deleted_at ON talent_trees(deleted_at);
CREATE INDEX IF NOT EXISTS idx_talent_nodes_deleted_at ON talent_nodes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_talent_entries_deleted_at ON talent_entries(deleted_at);
CREATE INDEX IF NOT EXISTS idx_sub_tree_nodes_deleted_at ON sub_tree_nodes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_sub_tree_entries_deleted_at ON sub_tree_entries(deleted_at);
CREATE INDEX IF NOT EXISTS idx_hero_nodes_deleted_at ON hero_nodes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_hero_entries_deleted_at ON hero_entries(deleted_at);