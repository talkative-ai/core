ALTER TABLE workbench_projects ADD COLUMN IF NOT EXISTS "Category" TEXT DEFAULT 0;
ALTER TABLE workbench_projects ADD COLUMN IF NOT EXISTS "Tags" TEXT[] DEFAULT ARRAY[]::TEXT[];

CREATE TABLE IF NOT EXISTS workbench_projects_needing_review (
    "ProjectID" UUID NOT NULL REFERENCES workbench_projects("ID"),
    "CreatedAt" timestamp DEFAULT current_timestamp
);
