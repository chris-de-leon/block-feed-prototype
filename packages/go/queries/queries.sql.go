// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package queries

import (
	"context"
)

const WebhooksFindOne = `-- name: WebhooksFindOne :one
SELECT id, created_at, is_active, url, max_blocks, max_retries, timeout_ms, customer_id, blockchain_id, shard_id FROM ` + "`" + `webhook` + "`" + ` WHERE ` + "`" + `id` + "`" + ` = ? LIMIT 1
`

// WebhooksFindOne
//
//	SELECT id, created_at, is_active, url, max_blocks, max_retries, timeout_ms, customer_id, blockchain_id, shard_id FROM `webhook` WHERE `id` = ? LIMIT 1
func (q *Queries) WebhooksFindOne(ctx context.Context, id string) (*Webhook, error) {
	row := q.db.QueryRowContext(ctx, WebhooksFindOne, id)
	var i Webhook
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.IsActive,
		&i.Url,
		&i.MaxBlocks,
		&i.MaxRetries,
		&i.TimeoutMs,
		&i.CustomerID,
		&i.BlockchainID,
		&i.ShardID,
	)
	return &i, err
}
