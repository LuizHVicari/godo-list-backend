-- name: CreateComment :exec
INSERT INTO todo.item_comments (id, item_id, author_id, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetCommentByID :one
SELECT id, item_id, author_id, content, created_at, updated_at FROM todo.item_comments WHERE id = $1;

-- name: CountCommentsByItemID :one
SELECT COUNT(*) FROM todo.item_comments WHERE item_id = $1;

-- name: ListCommentsByItemID :many
SELECT id, item_id, author_id, content, created_at, updated_at FROM todo.item_comments
WHERE item_id = $1
ORDER BY created_at ASC
LIMIT $3::INT
OFFSET $2::INT;

-- name: UpdateCommentByID :exec
UPDATE todo.item_comments SET content = $2, updated_at = $3 WHERE id = $1;

-- name: DeleteCommentByID :exec
DELETE FROM todo.item_comments WHERE id = $1;

-- name: IsItemInOwnedProject :one
SELECT EXISTS(
    SELECT 1 FROM todo.items i
    JOIN todo.steps s ON s.id = i.step_id
    JOIN todo.projects p ON p.id = s.project_id
    WHERE i.id = $1 AND p.owner_id = $2
) AS owned;
