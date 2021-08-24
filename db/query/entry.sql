-- name: CreateEntry :one
INSERT INTO entries (
          account_id, 
          amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE code = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY code
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :one
UPDATE  entries 
SET amount = $2
WHERE code = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE code = $1;