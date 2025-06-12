-- 022_fix_foreign_keys.up.sql
BEGIN;

-- 1. Correction des contraintes de la migration 003
-- Supprimer les contraintes existantes pour éviter les conflits
ALTER TABLE talent_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees;
ALTER TABLE talent_entries DROP CONSTRAINT IF EXISTS fk_talent_trees_entries;
ALTER TABLE sub_tree_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees_sub;
ALTER TABLE sub_tree_entries DROP CONSTRAINT IF EXISTS fk_sub_tree_nodes_entries;
ALTER TABLE hero_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees_hero;
ALTER TABLE hero_entries DROP CONSTRAINT IF EXISTS fk_hero_nodes_entries;

-- Recréer les contraintes proprement
ALTER TABLE talent_nodes
    ADD CONSTRAINT fk_talent_trees 
    FOREIGN KEY (talent_tree_id, spec_id) 
    REFERENCES talent_trees(trait_tree_id, spec_id);

ALTER TABLE talent_entries
    ADD CONSTRAINT fk_talent_trees_entries 
    FOREIGN KEY (talent_tree_id, spec_id) 
    REFERENCES talent_trees(trait_tree_id, spec_id);

ALTER TABLE sub_tree_nodes
    ADD CONSTRAINT fk_talent_trees_sub 
    FOREIGN KEY (talent_tree_id, spec_id) 
    REFERENCES talent_trees(trait_tree_id, spec_id);

ALTER TABLE sub_tree_entries
    ADD CONSTRAINT fk_sub_tree_nodes_entries 
    FOREIGN KEY (sub_tree_node_id) 
    REFERENCES sub_tree_nodes(sub_tree_node_id);

ALTER TABLE hero_nodes
    ADD CONSTRAINT fk_talent_trees_hero 
    FOREIGN KEY (talent_tree_id, spec_id) 
    REFERENCES talent_trees(trait_tree_id, spec_id);

ALTER TABLE hero_entries
    ADD CONSTRAINT fk_hero_nodes_entries 
    FOREIGN KEY (node_id, talent_tree_id, spec_id) 
    REFERENCES hero_nodes(node_id, talent_tree_id, spec_id);

-- 2. Correction pour la table dungeons et la clé étrangère
-- D'abord, vérifier si la table dungeons existe
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_name = 'dungeons'
    ) THEN
        -- Ajouter une contrainte UNIQUE à la colonne encounter_id de la table dungeons
        -- s'il n'existe pas déjà une telle contrainte
        IF NOT EXISTS (
            SELECT 1 
            FROM pg_constraint c
            JOIN pg_namespace n ON n.oid = c.connamespace
            WHERE c.contype = 'u' 
            AND c.conrelid = 'dungeons'::regclass
            AND array_position(c.conkey, (
                SELECT a.attnum 
                FROM pg_attribute a 
                WHERE a.attrelid = 'dungeons'::regclass 
                AND a.attname = 'encounter_id'
            )) IS NOT NULL
        ) THEN
            ALTER TABLE dungeons ADD CONSTRAINT uk_dungeons_encounter_id UNIQUE (encounter_id);
            RAISE NOTICE 'Added UNIQUE constraint to dungeons.encounter_id';
        ELSE
            RAISE NOTICE 'UNIQUE constraint already exists on dungeons.encounter_id';
        END IF;

        -- Maintenant on peut ajouter la clé étrangère
        ALTER TABLE warcraft_logs_reports
            DROP CONSTRAINT IF EXISTS fk_warcraft_logs_reports_encounter_id;
        
        -- Ne pas ajouter la contrainte si certaines valeurs encounter_id ne correspondent pas
        IF NOT EXISTS (
            SELECT wlr.encounter_id
            FROM warcraft_logs_reports wlr
            LEFT JOIN dungeons d ON wlr.encounter_id = d.encounter_id
            WHERE d.encounter_id IS NULL AND wlr.encounter_id IS NOT NULL
        ) THEN
            ALTER TABLE warcraft_logs_reports
                ADD CONSTRAINT fk_warcraft_logs_reports_encounter_id 
                FOREIGN KEY (encounter_id) 
                REFERENCES dungeons(encounter_id);
            RAISE NOTICE 'Foreign key constraint added successfully for warcraft_logs_reports';
        ELSE
            RAISE NOTICE 'Foreign key constraint not added because some encounter_ids do not exist in dungeons table';
        END IF;
    ELSE
        RAISE NOTICE 'Table dungeons does not exist, foreign key constraint not added';
    END IF;
END $$;

-- Ajoutons les colonnes et changements simples pour la migration 009
DO $$
BEGIN
    -- Ajout des colonnes
    ALTER TABLE warcraft_logs_reports 
        ADD COLUMN IF NOT EXISTS keystone_level INTEGER,
        ADD COLUMN IF NOT EXISTS keystone_time BIGINT,
        ADD COLUMN IF NOT EXISTS affixes INTEGER[],
        ADD COLUMN IF NOT EXISTS damage_taken JSONB;

    -- Modification des types de données
    -- Note: ces opérations peuvent échouer si les valeurs existantes ne sont pas convertibles
    BEGIN
        ALTER TABLE warcraft_logs_reports
            ALTER COLUMN fight_id SET DATA TYPE INTEGER,
            ALTER COLUMN encounter_id SET DATA TYPE INTEGER,
            ALTER COLUMN total_time SET DATA TYPE BIGINT,
            ALTER COLUMN item_level SET DATA TYPE NUMERIC;
        EXCEPTION WHEN OTHERS THEN
            RAISE NOTICE 'Could not convert column types: %', SQLERRM;
    END;

    -- Contraintes NOT NULL
    BEGIN
        ALTER TABLE warcraft_logs_reports
            ALTER COLUMN fight_id SET NOT NULL,
            ALTER COLUMN code SET NOT NULL;
        EXCEPTION WHEN OTHERS THEN
            RAISE NOTICE 'Could not set NOT NULL constraints: %', SQLERRM;
    END;
END $$;

-- Gestion des index
DROP INDEX IF EXISTS idx_reports_keystone;
CREATE INDEX IF NOT EXISTS idx_reports_encounter_id ON warcraft_logs_reports(encounter_id);
CREATE INDEX IF NOT EXISTS idx_reports_keystone_level ON warcraft_logs_reports(keystone_level);
CREATE INDEX IF NOT EXISTS idx_reports_keystone_time ON warcraft_logs_reports(keystone_time);

-- Index sur les données JSON seulement si la colonne existe
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'warcraft_logs_reports' AND column_name = 'damage_taken'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_reports_damage_taken ON warcraft_logs_reports USING gin (damage_taken jsonb_path_ops)';
    END IF;
END $$;

COMMIT;