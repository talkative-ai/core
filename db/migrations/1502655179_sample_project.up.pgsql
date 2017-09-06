INSERT INTO users ("Email", "GivenName", "FamilyName") VALUES ('warent@aum.ai', 'Wyatt', 'Arent');
INSERT INTO users ("Email", "GivenName", "FamilyName") VALUES ('cflack@aum.ai', 'Callum', 'Flack');
INSERT INTO teams ("Name") VALUES (NULL);
INSERT INTO team_members ("TeamID", "UserID", "Role") VALUES (1, 1, 1);
INSERT INTO team_members ("TeamID", "UserID", "Role") VALUES (1, 2, 1);

INSERT INTO workbench_projects ("TeamID", "Title") VALUES (1, 'The Spooky House');
INSERT INTO workbench_zones ("ProjectID", "Title") VALUES (1, 'Hall');
INSERT INTO workbench_actors ("ProjectID", "Title") VALUES (1, 'Kilroy');
INSERT INTO workbench_zones_actors ("ZoneID", "ActorID") VALUES (1, 1);
UPDATE workbench_projects SET "StartZoneID"=1 WHERE "ID"=1;
INSERT INTO workbench_dialog_nodes ("ActorID", "EntryInput", "AlwaysExec", "IsRoot") VALUES (1, '{Greetings}', '{"PlaySounds":[{"SoundType":0,"Val":"Hello world!"}]}', true);
INSERT INTO workbench_dialog_nodes ("ActorID", "EntryInput", "AlwaysExec") VALUES (1, '{Farewell}', '{"PlaySounds":[{"SoundType":0,"Val":"It was nice talking to you."}]}');
INSERT INTO workbench_dialog_nodes_relations ("ParentNodeID", "ChildNodeID") VALUES (1, 2);