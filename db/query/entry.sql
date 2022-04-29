-- name: CreateEntry :execresult
INSERT INTO entries(
    `account_id`, `amount`
) VALUES (?, ?);

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ? LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateEntry :exec
UPDATE entries
SET amount = ?
WHERE id = ?;

-- name: DeleteEntry :exec
DELETE
FROM entries
WHERE id = ?;