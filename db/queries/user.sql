-- name: CreateUser :exec
INSERT INTO auth.users (
    id,
    email,
    password_hash,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM auth.users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM auth.users WHERE id = $1;

-- name: UpdateUserByID :exec
UPDATE auth.users SET
    email = $2,
    password_hash = $3,
    updated_at = $4
WHERE id = $1 RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM auth.users WHERE id = $1;