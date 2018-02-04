ALTER TABLE static_published_projects_versioned ALTER "ProjectData" SET NOT NULL;
ALTER TABLE static_published_projects_versioned ALTER "ProjectData" SET DEFAULT '[]'::jsonb;
ALTER TABLE static_published_projects_versioned ADD COLUMN IF NOT EXISTS "TriggerData" JSONB NOT NULL DEFAULT '[]'::jsonb;