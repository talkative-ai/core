CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users
CREATE TABLE IF NOT EXISTS users (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "GivenName" TEXT,
    "FamilyName" TEXT,
    "Email" TEXT,
    "Image" TEXT,
    "PasswordSHA" TEXT,
    "Salt" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS upgrade_item (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "UserID" UUID NOT NULL REFERENCES users("ID"),
    "SKU" INTEGER NOT NULL,
    "Trial" timestamp
);

-- User Teams
CREATE TABLE IF NOT EXISTS teams (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "Name" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS team_members (
    "TeamID" UUID NOT NULL REFERENCES teams("ID"),
    "UserID" UUID NOT NULL REFERENCES users("ID"),
    "Role" SMALLINT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

-- Workbench Projects
CREATE TABLE IF NOT EXISTS workbench_projects (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "TeamID" UUID NOT NULL REFERENCES teams("ID"),
    "Title" TEXT,
    "StartZoneID" UUID,
    "CreatedAt" timestamp DEFAULT current_timestamp,
    "IsPrivate" BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS workbench_private_project_grants (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "ProjectID" UUID NOT NULL REFERENCES workbench_projects("ID"),
    "UserID" UUID NOT NULL REFERENCES users("ID")
);

-- Basic entities
CREATE TABLE IF NOT EXISTS workbench_actors (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "ProjectID" UUID NOT NULL REFERENCES workbench_projects("ID"),
    "Title" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS workbench_zones (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "ProjectID" UUID NOT NULL REFERENCES workbench_projects("ID"),
    "Title" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

ALTER TABLE workbench_projects
    ADD CONSTRAINT fk_startzone
    FOREIGN KEY ("StartZoneID")
    REFERENCES workbench_zones("ID");

CREATE TABLE IF NOT EXISTS workbench_triggers (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "TriggerType" INTEGER,
    "ProjectID" UUID NOT NULL REFERENCES workbench_projects("ID"),
    "ZoneID" UUID NOT NULL REFERENCES workbench_zones("ID"),
    "AlwaysExec" JSONB,
    "Statements" JSONB,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS workbench_dialog_nodes (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "ActorID" UUID NOT NULL REFERENCES workbench_actors("ID"),
    "EntryInput" TEXT[],
    "AlwaysExec" JSONB,
    "Statements" JSONB,
    "IsRoot" BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS workbench_dialog_nodes_relations (
    "ParentNodeID" UUID NOT NULL REFERENCES workbench_dialog_nodes("ID"),
    "ChildNodeID" UUID NOT NULL REFERENCES workbench_dialog_nodes("ID")
);

CREATE TABLE IF NOT EXISTS workbench_zones_actors (
    "ZoneID" UUID NOT NULL REFERENCES workbench_zones("ID"),
    "ActorID" UUID NOT NULL REFERENCES workbench_actors("ID")
);

-- Gameplay event sourcing tables
CREATE TABLE IF NOT EXISTS event_user_action (
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "UserID" UUID NOT NULL REFERENCES users("ID"),
    "PubID" UUID NOT NULL REFERENCES workbench_projects("ID"), -- Publishing ID; a unique ID that the project will have when published
    "RawInput" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS event_state_change (
    "EventUserActionID" UUID NOT NULL REFERENCES event_user_action("ID"),
    "StateObject" JSONB,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

-- Misc
CREATE TABLE IF NOT EXISTS published_workbench_projects (
    "ProjectID" UUID NOT NULL REFERENCES workbench_projects("ID"),
    "TeamID" UUID NOT NULL REFERENCES teams("ID"),
    "CreatedAt" timestamp DEFAULT current_timestamp
);