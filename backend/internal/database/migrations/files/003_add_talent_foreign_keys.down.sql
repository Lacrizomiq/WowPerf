-- 003_add_talent_foreign_keys.down.sql
ALTER TABLE hero_entries DROP CONSTRAINT IF EXISTS fk_hero_nodes_entries;
ALTER TABLE hero_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees_hero;
ALTER TABLE sub_tree_entries DROP CONSTRAINT IF EXISTS fk_sub_tree_nodes_entries;
ALTER TABLE sub_tree_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees_sub;
ALTER TABLE talent_entries DROP CONSTRAINT IF EXISTS fk_talent_trees_entries;
ALTER TABLE talent_nodes DROP CONSTRAINT IF EXISTS fk_talent_trees;