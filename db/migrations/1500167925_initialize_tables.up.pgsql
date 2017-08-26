-- Users
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT,
    passwordsha TEXT,
    salt TEXT
);

CREATE TABLE IF NOT EXISTS user_linked_accounts (
    user_id BIGINT NOT NULL REFERENCES users(id),
    provider TEXT,
    email TEXT
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

-- Projects
CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL REFERENCES teams(id),
    title TEXT,
    start_zone_id BIGINT,
    created_at timestamp DEFAULT current_timestamp
);

-- Basic entities
CREATE TABLE IF NOT EXISTS actors (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id),
    title TEXT,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS zones (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id),
    title TEXT,
    created_at timestamp DEFAULT current_timestamp
);

ALTER TABLE projects
    ADD CONSTRAINT fk_startzone
    FOREIGN KEY (start_zone_id)
    REFERENCES zones(id);

CREATE TABLE IF NOT EXISTS triggers (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id),
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS notes (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id),
    title TEXT,
    content TEXT,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS logical_set (
    id BIGSERIAL PRIMARY KEY,
    always JSONB,
    statements JSONB
);
CREATE TABLE IF NOT EXISTS dialog_nodes (
    id BIGSERIAL PRIMARY KEY,
    zone_id BIGINT NOT NULL REFERENCES zones(id),
    entry TEXT[],
    logical_set_id BIGINT NOT NULL REFERENCES logical_set(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS dialog_nodes_relations (
    parent_node_id BIGINT NOT NULL REFERENCES dialog_nodes(id),
    child_node_id BIGINT NOT NULL REFERENCES dialog_nodes(id)
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
CREATE TABLE IF NOT EXISTS published_projects (
    pub_id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id),
    title TEXT,
    creator TEXT,
    created_at timestamp DEFAULT current_timestamp

);