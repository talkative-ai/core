INSERT INTO projects (title) VALUES ('The Spooky House');
INSERT INTO zones (project_id, title) VALUES (1, 'Hall');
UPDATE projects SET start_zone_id=1 WHERE id=1;
INSERT INTO logical_set (always) VALUES ('{"play_sounds":[{"sound_type":0,"value":"Hello world!"}]}');
INSERT INTO logical_set (always) VALUES ('{"play_sounds":[{"sound_type":0,"value":"It was nice talking to you."}]}');
INSERT INTO dialog_nodes (zone_id, entry, logical_set_id) VALUES (1, '{"statement_greeting"}', 1);
INSERT INTO dialog_nodes (zone_id, entry, logical_set_id) VALUES (1, '{"statement_farewell"}', 2);
INSERT INTO dialog_nodes_relations (parent_node_id, child_node_id) VALUES (1, 2);