-- Users
CREATE TABLE IF NOT EXISTS users (
    "ID" BIGSERIAL PRIMARY KEY,
    "GivenName" TEXT,
    "FamilyName" TEXT,
    "Email" TEXT,
    "Image" TEXT,
    "PasswordSHA" TEXT,
    "Salt" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

-- User Teams
CREATE TABLE IF NOT EXISTS teams (
    "ID" BIGSERIAL PRIMARY KEY,
    "Name" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS team_members (
    "TeamID" BIGINT NOT NULL REFERENCES teams("ID"),
    "UserID" BIGINT NOT NULL REFERENCES users("ID"),
    "Role" SMALLINT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

-- Workbench Projects
CREATE TABLE IF NOT EXISTS workbench_projects (
    "ID" BIGSERIAL PRIMARY KEY,
    "TeamID" BIGINT NOT NULL REFERENCES teams("ID"),
    "Title" TEXT,
    "StartZoneID" BIGINT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

-- Basic entities
CREATE TABLE IF NOT EXISTS workbench_actors (
    "ID" BIGSERIAL PRIMARY KEY,
    "ProjectID" BIGINT NOT NULL REFERENCES workbench_projects("ID"),
    "Title" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS workbench_zones (
    "ID" BIGSERIAL PRIMARY KEY,
    "ProjectID" BIGINT NOT NULL REFERENCES workbench_projects("ID"),
    "Title" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

ALTER TABLE workbench_projects
    ADD CONSTRAINT fk_startzone
    FOREIGN KEY ("StartZoneID")
    REFERENCES workbench_zones("ID");

CREATE TABLE IF NOT EXISTS workbench_triggers (
    "ID" BIGSERIAL PRIMARY KEY,
    "ProjectID" BIGINT NOT NULL REFERENCES workbench_projects("ID"),
    "CreatedAt" timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS workbench_notes (
    "ID" BIGSERIAL PRIMARY KEY,
    "ProjectID" BIGINT NOT NULL REFERENCES workbench_projects("ID"),
    "Title" TEXT,
    "Content" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS workbench_logical_set (
    "ID" BIGSERIAL PRIMARY KEY,
    "Always" JSONB,
    "Statements" JSONB
);
CREATE TABLE IF NOT EXISTS workbench_dialog_nodes (
    "ID" BIGSERIAL PRIMARY KEY,
    "ActorID" BIGINT NOT NULL REFERENCES workbench_actors("ID"),
    "Entry" TEXT[],
    "LogicalSetID" BIGINT NOT NULL REFERENCES workbench_logical_set("ID")
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workbench_dialog_nodes_relations (
    "ParentNodeID" BIGINT NOT NULL REFERENCES workbench_dialog_nodes("ID"),
    "ChildNodeID" BIGINT NOT NULL REFERENCES workbench_dialog_nodes("ID")
);

CREATE TABLE IF NOT EXISTS workbench_zones_actors (
    "ID" BIGSERIAL PRIMARY KEY,
    "ZoneID" BIGINT NOT NULL REFERENCES workbench_zones("ID"),
    "ActorID" BIGINT NOT NULL REFERENCES workbench_actors("ID")
);

-- Gameplay event sourcing tables
CREATE TABLE IF NOT EXISTS event_user_action (
    "ID" BIGSERIAL PRIMARY KEY,
    "UserID" BIGINT NOT NULL REFERENCES users("ID"),
    "PubID" integer, -- Publishing ID; a unique ID that the project will have when published
    "RawInput" TEXT,
    "CreatedAt" timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS event_state_change (
    "EventUserActionID" BIGINT NOT NULL REFERENCES event_user_action("ID"),
    "StateObject" JSONB,
    "CreatedAt" timestamp DEFAULT current_timestamp
);

-- Misc
CREATE TABLE IF NOT EXISTS published_workbench_projects (
    "PubID" BIGSERIAL PRIMARY KEY,
    "ProjectID" BIGINT NOT NULL REFERENCES workbench_projects("ID"),
    "Title" TEXT,
    "Publisher" BIGINT NOT NULL REFERENCES teams("ID"),
    "CreatedAt" timestamp DEFAULT current_timestamp

);