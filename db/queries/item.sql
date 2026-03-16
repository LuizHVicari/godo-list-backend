-- name: CreateItem :exec
INSERT INTO todo.items
    (id, step_id, name, description, priority, position, is_completed, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetItemByID :one
SELECT * FROM todo.items WHERE id = $1;

-- name: CountItemsByStepID :one
SELECT COUNT(*) FROM todo.items
WHERE step_id = $1
  AND (sqlc.narg('name')::TEXT IS NULL
  OR name ILIKE '%' || sqlc.narg('name')::TEXT || '%');

-- name: ListItemsByStepID :many
SELECT * FROM todo.items
WHERE step_id = $1
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
    CASE WHEN sqlc.narg('sort')::TEXT = 'priority'
            AND sqlc.narg('direction')::TEXT = 'asc'  THEN priority::TEXT END ASC,
    CASE WHEN sqlc.narg('sort')::TEXT = 'priority'
            AND (sqlc.narg('direction')::TEXT = 'desc'
            OR sqlc.narg('direction')::TEXT IS NULL)   THEN priority::TEXT END DESC,
    CASE WHEN (sqlc.narg('sort')::TEXT = 'position' OR sqlc.narg('sort')::TEXT IS NULL)
            AND sqlc.narg('direction')::TEXT = 'desc' THEN position END DESC,
    CASE WHEN (sqlc.narg('sort')::TEXT = 'position' OR sqlc.narg('sort')::TEXT IS NULL)
            AND (sqlc.narg('direction')::TEXT = 'asc'
            OR sqlc.narg('direction')::TEXT IS NULL)   THEN position END ASC
LIMIT sqlc.narg('limit')::INT
OFFSET sqlc.narg('offset')::INT;

-- name: IsStepInOwnedProject :one
SELECT EXISTS(
    SELECT 1 FROM todo.steps s
    JOIN todo.projects p ON p.id = s.project_id
    WHERE s.id = $1 AND p.owner_id = $2
) AS owned;

-- name: GetLastItemPositionByStepID :one
SELECT COALESCE(MAX(position), 0)::INT FROM todo.items WHERE step_id = $1;

-- name: IsItemPositionTaken :one
SELECT EXISTS(
    SELECT 1 FROM todo.items
    WHERE step_id = $1 AND position = $2
      AND (sqlc.narg('exclude_id')::UUID IS NULL OR id != sqlc.narg('exclude_id')::UUID)
) AS taken;

-- name: UpdateItemPositionByID :execresult
UPDATE todo.items SET position = $3, updated_at = $4 WHERE id = $1 AND step_id = $2;

-- name: UpdateItemByID :exec
UPDATE todo.items SET
    step_id = $2,
    name = $3,
    description = $4,
    priority = $5,
    position = $6,
    is_completed = $7,
    updated_at = $8
WHERE id = $1;

-- name: DeleteItemByID :exec
DELETE FROM todo.items WHERE id = $1;
