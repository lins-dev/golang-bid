-- name: CreateProduct :one
INSERT INTO products("seller_id", "product_name", "description", "price", "auction_end")
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1
LIMIT 1;

-- name: GetProductByUuid :one
SELECT * FROM products
WHERE uuid = $1
LIMIT 1;

-- name: ListProduct :many
SELECT * FROM products
ORDER BY product_name;

-- name: ListProductPaginated :many
SELECT * FROM products
ORDER BY product_name
LIMIT $1 OFFSET $2;