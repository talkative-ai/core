-- Basic entities
CREATE TABLE IF NOT EXISTS actors (
    id BIGSERIAL PRIMARY KEY,
    title text,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS dialogs (
    id BIGSERIAL PRIMARY KEY,
    title text,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS zones (
    id BIGSERIAL PRIMARY KEY,
    title text,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS triggers (
    id BIGSERIAL PRIMARY KEY,
    created_at
);
CREATE TABLE IF NOT EXISTS notes (
    id BIGSERIAL PRIMARY KEY,
    title text,
    content text,
    created_at timestamp DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS action_sets (
    id BIGSERIAL PRIMARY KEY,
    always jsonb,
    statements jsonb
);

-- Projects
CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    title text,
    owner_id text,
    start_zone_id integer references zones(id),
    created_at timestamp DEFAULT current_timestamp
);

-- Relate entities to project
CREATE TABLE IF NOT EXISTS project_actors (
    project_id integer NOT NULL references projects(id),
    actor_id integer NOT NULL references actors(id)
);
CREATE TABLE IF NOT EXISTS project_dialogs (
    project_id integer NOT NULL references projects(id),
    dialog_id integer NOT NULL references dialogs(id)
);
CREATE TABLE IF NOT EXISTS project_zones (
    project_id integer NOT NULL references projects(id),
    zone_id integer NOT NULL references zones(id)
);
CREATE TABLE IF NOT EXISTS project_notes (
    project_id integer NOT NULL references projects(id),
    note_id integer NOT NULL references notes(id)
);

-- Relate triggers to zones
CREATE TABLE IF NOT EXISTS zone_triggers (
    zone_id integer NOT NULL references zones(id)
    trigger_id integer NOT NULL references triggers(id)
);

-- Relate action sets
CREATE TABLE IF NOT EXISTS trigger_action_sets (
    trigger_id integer NOT NULL references triggers(id),
    action_set_id integer NOT NULL references action_sets(id)
);
CREATE TABLE IF NOT EXISTS dialog_node_action_sets (
    dialog_node_id integer NOT NULL references dialog_nodes(id),
    action_set_id integer NOT NULL references action_sets(id)
);


-- Gameplay event sourcing tables
CREATE TABLE IF NOT EXISTS event_user_action (
    id BIGSERIAL PRIMARY KEY,
    user_id integer,
    pub_id integer, -- Publishing ID; a unique ID that the project will have when published
    raw_input string,
    created_at
);
CREATE TABLE IF NOT EXISTS event_state_change (
    event_user_action_id integer NOT NULL references event_user_action(id),
    state_object jsonb,
    created_at timestamp DEFAULT current_timestamp
);


CREATE TABLE IF NOT EXISTS published_projects (
    title text,
    pub_id PRIMARY KEY,
    project_id integer,
    creator text
)