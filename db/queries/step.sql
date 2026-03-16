-- name: CreateStep :exec
INSERT INTO todo.steps
    (id, project_id, name, position, is_completed, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7);

-- name: GetStepByID :one
SELECT * FROM todo.steps WHERE id = $1;

-- name: CountStepsByProjectID :one
SELECT COUNT(*) FROM todo.steps
WHERE project_id = $1
  AND (sqlc.narg('name')::TEXT IS NULL
  OR name ILIKE '%' || sqlc.narg('name')::TEXT || '%');

-- name: ListStepsByProjectID :many
SELECT * FROM todo.steps
WHERE project_id = $1
  AND (sqlc.narg('name')::TEXT IS NULL
  OR name ILIKE '%' || sqlc.narg('name')::TEXT || '%')
ORDER BY
    CASE WHEN sqlc.narg('sort')::TEXT = 'name'
            AND sqlc.narg('direction')::TEXT = 'asc'  THEN name END ASC,
    CASE WHEN sqlc.narg('sort')::TEXT = 'name'
            AND (sqlc.narg('direction')::TEXT = 'desc'
            OR sqlc.narg('direction')::TEXT IS NULL)   THEN name END DESC,
    CASE WHEN sqlc.narg('sort')::TEXT = 'updated_at'
            AND sqlc.narg('direction')::TEXT = 'asc'  THEN updated_at END ASC,
    CASE WHEN sqlc.narg('sort')::TEXT = 'updated_at'
            AND (sqlc.narg('direction')::TEXT = 'desc'
            OR sqlc.narg('direction')::TEXT IS NULL)   THEN updated_at END DESC,
    CASE WHEN sqlc.narg('sort')::TEXT = 'created_at'
            AND sqlc.narg('direction')::TEXT = 'asc'  THEN created_at END ASC,
    CASE WHEN sqlc.narg('sort')::TEXT = 'created_at'
            AND (sqlc.narg('direction')::TEXT = 'desc'
            OR sqlc.narg('direction')::TEXT IS NULL)   THEN created_at END DESC,
    CASE WHEN (sqlc.narg('sort')::TEXT = 'position' OR sqlc.narg('sort')::TEXT IS NULL)
            AND sqlc.narg('direction')::TEXT = 'desc' THEN position END DESC,
    CASE WHEN (sqlc.narg('sort')::TEXT = 'position' OR sqlc.narg('sort')::TEXT IS NULL)
            AND (sqlc.narg('direction')::TEXT = 'asc'
            OR sqlc.narg('direction')::TEXT IS NULL)   THEN position END ASC
LIMIT sqlc.narg('limit')::INT
OFFSET sqlc.narg('offset')::INT;

-- name: GetLastStepPositionByProjectID :one
SELECT COALESCE(MAX(position), 0)::INT FROM todo.steps WHERE project_id = $1;

-- name: UpdateStepPositionByID :execresult
UPDATE todo.steps SET position = $3, updated_at = $4 WHERE id = $1 AND project_id = $2;

-- name: UpdateStepByID :exec
UPDATE todo.steps SET
    name = $2,
    position = $3,
    is_completed = $4,
    updated_at = $5
WHERE id = $1;

-- name: DeleteStepByID :exec
DELETE FROM todo.steps WHERE id = $1;
