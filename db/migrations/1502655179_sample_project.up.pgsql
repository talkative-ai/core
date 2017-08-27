-- Password: THISisthePASSWORD!forAUMtestACCOUNT:)
INSERT INTO users ("Email", "PasswordSHA", "Salt") VALUES (
  'test@aum.ai',
  '369214D341838EADED50CB089766B8EDBB74F9AB642059EAA7CD5F977F80B57ED8F455C0EC9A13979CEEB5DB2D65FF121AE2E817FDE285BFE548D9128177FC27',
  'qrbFIauvCuipVogmzvJB');
INSERT INTO teams ("Name") VALUES (NULL);
INSERT INTO team_members ("TeamID", "UserID", "Role") VALUES (1, 1, 1);

INSERT INTO workbench_projects ("TeamID", "Title") VALUES (1, 'The Spooky House');
INSERT INTO workbench_zones ("ProjectID", "Title") VALUES (1, 'Hall');
UPDATE workbench_projects SET "StartZoneID"=1 WHERE "ID"=1;
INSERT INTO workbench_logical_set ("Always") VALUES ('{"play_sounds":[{"sound_type":0,"value":"Hello world!"}]}');
INSERT INTO workbench_logical_set ("Always") VALUES ('{"play_sounds":[{"sound_type":0,"value":"It was nice talking to you."}]}');
INSERT INTO workbench_dialog_nodes ("ZoneID", "Entry", "LogicalSetID") VALUES (1, '{"statement_greeting"}', 1);
INSERT INTO workbench_dialog_nodes ("ZoneID", "Entry", "LogicalSetID") VALUES (1, '{"statement_farewell"}', 2);
INSERT INTO workbench_dialog_nodes_relations ("ParentNodeID", "ChildNodeID") VALUES (1, 2);