-- +goose Up
CREATE TYPE todo.item_priority AS ENUM ('none', 'low', 'medium', 'high', 'critical');

CREATE TABLE todo.items (
    id UUID PRIMARY KEY,
    step_id UUID NOT NULL REFERENCES todo.steps(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    priority todo.item_priority NOT NULL DEFAULT 'none',
    position INT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (step_id, position) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_items_step_id ON todo.items(step_id);

-- +goose Down
DROP TABLE IF EXISTS todo.items;
DROP TYPE IF EXISTS todo.item_priority;
