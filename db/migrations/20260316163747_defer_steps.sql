-- +goose Up
ALTER TABLE todo.steps DROP CONSTRAINT IF EXISTS steps_project_id_position_key;
ALTER TABLE todo.steps ADD CONSTRAINT steps_project_id_position_key UNIQUE (project_id, position) DEFERRABLE INITIALLY DEFERRED;

-- +goose Down
ALTER TABLE todo.steps DROP CONSTRAINT IF EXISTS steps_project_id_position_key;
ALTER TABLE todo.steps ADD CONSTRAINT steps_project_id_position_key UNIQUE (project_id, position);
