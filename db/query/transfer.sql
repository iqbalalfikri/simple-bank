-- name: CreateTransfer :execresult
INSERT IGNORE INTO transfers(
    `from_account_id`, `to_account_id`, `amount`
) VALUES (?, ?, ?);

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ? LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateTransfer :exec
UPDATE transfers
SET amount = ?
WHERE id = ?;

-- name: DeleteTransfer :exec
DELETE
FROM transfers
WHERE id = ?;