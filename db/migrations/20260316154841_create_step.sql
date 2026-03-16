-- +goose Up
CREATE TABLE todo.steps (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES todo.projects(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name TEXT NOT NULL,
    position INT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (project_id, position)
);

CREATE INDEX idx_steps_project_id ON todo.steps(project_id);

-- +goose Down
DROP TABLE IF EXISTS todo.steps;