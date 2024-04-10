-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    card_number
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;
