-- name: CreateProject :exec
INSERT INTO todo.projects 
    (id, name, description, owner_id, created_at, updated_at) 
VALUES 
    ($1, $2, $3, $4, $5, $6);

-- name: GetProjectByID :one
SELECT * FROM todo.projects WHERE id = $1;

-- name: ListProjectsByOwnerID :many
SELECT * FROM todo.projects
WHERE owner_id = $1
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
    CASE WHEN (sqlc.narg('sort')::TEXT = 'created_at' OR sqlc.narg('sort')::TEXT IS NULL) 
            AND sqlc.narg('direction')::TEXT = 'asc'  THEN created_at END ASC,
    CASE WHEN (sqlc.narg('sort')::TEXT = 'created_at' OR sqlc.narg('sort')::TEXT IS NULL) 
            AND (sqlc.narg('direction')::TEXT = 'desc' 
            OR sqlc.narg('direction')::TEXT IS NULL)   THEN created_at END DESC
LIMIT sqlc.narg('limit')::INT
OFFSET sqlc.narg('offset')::INT;

-- name: CountProjectsByOwnerID :one
SELECT COUNT(*) FROM todo.projects
WHERE owner_id = $1
  AND (sqlc.narg('name')::TEXT IS NULL
  OR name ILIKE '%' || sqlc.narg('name')::TEXT || '%');

-- name: UpdateProjectByID :exec
UPDATE todo.projects SET
    name = $2,
    description = $3,
    updated_at = $4
WHERE id = $1;

-- name: DeleteProjectByID :exec
DELETE FROM todo.projects WHERE id = $1;