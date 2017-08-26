-- Users
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    email TEXT,
    image TEXT,
    passwordsha TEXT,
    salt TEXT
);

-- User Teams
CREATE TABLE IF NOT EXISTS teams (
    id BIGSERIAL PRIMARY KEY,
    team_name TEXT
);

CREATE TABLE IF NOT EXISTS team_members (
    team_id BIGINT NOT NULL REFERENCES teams(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    role SMALLINT
);

-- Workbench Projects
CREATE TABLE IF NOT EXISTS workbench_projects (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL REFERENCES teams(id),
    title TEXT,
    start_zone_id BIGINT,
    created_at timestamp DEFAULT current_timestamp
);

-- Basic entities
CREATE TABLE IF NOT EXISTS workbench_actors (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES workbench_projects(id),
    title TEXT,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS workbench_zones (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES workbench_projects(id),
    title TEXT,
    created_at timestamp DEFAULT current_timestamp
);

ALTER TABLE workbench_projects
    ADD CONSTRAINT fk_startzone
    FOREIGN KEY (start_zone_id)
    REFERENCES workbench_zones(id);

CREATE TABLE IF NOT EXISTS workbench_triggers (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES workbench_projects(id),
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS workbench_notes (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES workbench_projects(id),
    title TEXT,
    content TEXT,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS workbench_logical_set (
    id BIGSERIAL PRIMARY KEY,
    always JSONB,
    statements JSONB
);
CREATE TABLE IF NOT EXISTS workbench_dialog_nodes (
    id BIGSERIAL PRIMARY KEY,
    zone_id BIGINT NOT NULL REFERENCES workbench_zones(id),
    entry TEXT[],
    logical_set_id BIGINT NOT NULL REFERENCES workbench_logical_set(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workbench_dialog_nodes_relations (
    parent_node_id BIGINT NOT NULL REFERENCES workbench_dialog_nodes(id),
    child_node_id BIGINT NOT NULL REFERENCES workbench_dialog_nodes(id)
);

-- Gameplay event sourcing tables
CREATE TABLE IF NOT EXISTS event_user_action (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    pub_id integer, -- Publishing ID; a unique ID that the project will have when published
    raw_input TEXT,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS event_state_change (
    event_user_action_id BIGINT NOT NULL REFERENCES event_user_action(id),
    state_object JSONB,
    created_at timestamp DEFAULT current_timestamp
);

-- Misc
CREATE TABLE IF NOT EXISTS published_workbench_projects (
    pub_id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES workbench_projects(id),
    title TEXT,
    creator TEXT,
    created_at timestamp DEFAULT current_timestamp

);