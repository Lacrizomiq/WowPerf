-- 022_fix_foreign_keys.down.sql
BEGIN;

-- Supprimer les contraintes ajoutées/corrigées
ALTER TABLE talent_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees;
ALTER TABLE talent_entries DROP CONSTRAINT IF EXISTS fk_talent_trees_entries;
ALTER TABLE sub_tree_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees_sub;
ALTER TABLE sub_tree_entries DROP CONSTRAINT IF EXISTS fk_sub_tree_nodes_entries;
ALTER TABLE hero_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees_hero;
ALTER TABLE hero_entries DROP CONSTRAINT IF EXISTS fk_hero_nodes_entries;

-- Supprimer la contrainte de clé étrangère pour warcraft_logs_reports
ALTER TABLE warcraft_logs_reports
    DROP CONSTRAINT IF EXISTS fk_warcraft_logs_reports_encounter_id;

-- Supprimer la contrainte UNIQUE sur dungeons
ALTER TABLE dungeons DROP CONSTRAINT IF EXISTS uk_dungeons_encounter_id;

-- Supprimer les index créés
DROP INDEX IF EXISTS idx_reports_encounter_id;
DROP INDEX IF EXISTS idx_reports_keystone_level;
DROP INDEX IF EXISTS idx_reports_keystone_time;
DROP INDEX IF EXISTS idx_reports_damage_taken;

COMMIT;