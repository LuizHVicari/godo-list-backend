-- +goose Up
CREATE SCHEMA IF NOT EXISTS todo;
CREATE TABLE todo.projects (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_owner_id ON todo.projects(owner_id);
CREATE INDEX idx_projects_name ON todo.projects(name);


-- +goose Down
DROP TABLE IF EXISTS todo.projects;
DROP SCHEMA IF EXISTS todo;