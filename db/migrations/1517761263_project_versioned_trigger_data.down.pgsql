ALTER TABLE static_published_projects_versioned ALTER "ProjectData" DROP NOT NULL;
ALTER TABLE static_published_projects_versioned ALTER "ProjectData" DROP DEFAULT;
ALTER TABLE static_published_projects_versioned DROP COLUMN IF EXISTS "TriggerData";