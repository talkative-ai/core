-- Projects
CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    title text,
    owner_id text,
    start_zone_id bigint,
    created_at timestamp DEFAULT current_timestamp
);

-- Basic entities
CREATE TABLE IF NOT EXISTS actors (
    id BIGSERIAL PRIMARY KEY,
    project_id bigint NOT NULL references projects(id),
    title text,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS zones (
    id BIGSERIAL PRIMARY KEY,
    project_id bigint NOT NULL references projects(id),
    title text,
    created_at timestamp DEFAULT current_timestamp
);

ALTER TABLE projects
    ADD CONSTRAINT fk_startzone
    FOREIGN KEY (start_zone_id)
    REFERENCES zones(id);

CREATE TABLE IF NOT EXISTS triggers (
    id BIGSERIAL PRIMARY KEY,
    project_id bigint NOT NULL references projects(id),
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS notes (
    id BIGSERIAL PRIMARY KEY,
    project_id bigint NOT NULL references projects(id),
    title text,
    content text,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS logical_set (
    id BIGSERIAL PRIMARY KEY,
    always JSONB,
    statements JSONB
);
CREATE TABLE IF NOT EXISTS dialog_nodes (
    id BIGSERIAL PRIMARY KEY,
    zone_id bigint NOT NULL references zones(id),
    entry text[],
    logical_set_id bigint NOT NULL references logical_set(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS dialog_nodes_relations (
    parent_node_id bigint NOT NULL references dialog_nodes(id),
    child_node_id bigint NOT NULL references dialog_nodes(id)
);

-- Gameplay event sourcing tables
CREATE TABLE IF NOT EXISTS event_user_action (
    id BIGSERIAL PRIMARY KEY,
    user_id integer,
    pub_id integer, -- Publishing ID; a unique ID that the project will have when published
    raw_input text,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS event_state_change (
    event_user_action_id bigint NOT NULL references event_user_action(id),
    state_object JSONB,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS published_projects (
    pub_id BIGSERIAL PRIMARY KEY,
    project_id bigint NOT NULL references projects(id),
    title text,
    creator text
);