INSERT INTO projects (title) VALUES ('The Spooky House');
INSERT INTO zones (project_id, title) VALUES (1, 'Hall');
INSERT INTO logical_set (always) VALUES ('{"play_sounds":[{"sound_type":0,"value":"Hello world!"}]}');
INSERT INTO dialog_nodes (zone_id, entry, logical_set_id) VALUES (1, '{"statement_greeting"}', 1);
