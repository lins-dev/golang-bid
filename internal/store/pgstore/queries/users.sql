-- name: CreateUser :one
INSERT INTO users("user_name", "email", "password_hash", "bio")
VALUES($1, $2, $3, $4)
    RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUuid :one
SELECT * FROM users
WHERE uuid = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY email;

-- name: UpdateUser :one
UPDATE users
    set user_name = $2,
    email = $3,
    password_hash = $4,
    bio = $5
WHERE id = $1
RETURNING *;

-- name: UpdateUserByUuid :one
UPDATE users
    set user_name = $2,
    email = $3,
    password_hash = $4,
    bio = $5
WHERE uuid = $1
RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;