CREATE TABLE IF NOT EXISTS static_published_projects_versioned (
    "ProjectID" UUID,
    "Version" BIGINT,
    "Title" TEXT,
    "Category" TEXT,
    "Tags" TEXT[] DEFAULT ARRAY[]::TEXT[],
    "ProjectData" JSONB
);