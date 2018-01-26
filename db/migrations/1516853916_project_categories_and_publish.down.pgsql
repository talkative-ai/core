ALTER TABLE workbench_projects DROP COLUMN IF EXISTS "Category";
ALTER TABLE workbench_projects DROP COLUMN IF EXISTS "Tags";

DROP TABLE IF EXISTS workbench_projects_needing_review CASCADE;