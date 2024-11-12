-- 003_add_talent_foreign_keys.up.sql
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