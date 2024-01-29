// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package sqlc

import (
	"context"
)

type CreatePendingWebhookJobParams struct {
	BlockHeight string `json:"blockHeight"`
	ChainID     string `json:"chainId"`
	ChainUrl    string `json:"chainUrl"`
	ChannelName string `json:"channelName"`
}

const createWebhook = `-- name: CreateWebhook :execrows
WITH inserted_user AS (
  INSERT INTO "customer" ("id", "created_at")
  VALUES ($6::TEXT, DEFAULT)
  ON CONFLICT ("id") DO NOTHING
  RETURNING "id"
)
INSERT INTO "webhook" (
  "id",
  "chain_id",
  "url",
  "max_retries",
  "timeout_ms",
  "retry_delay_ms",
  "customer_id"
) VALUES (
  DEFAULT,
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
`

type CreateWebhookParams struct {
	ChainID      string `json:"chainId"`
	Url          string `json:"url"`
	MaxRetries   int32  `json:"maxRetries"`
	TimeoutMs    int32  `json:"timeoutMs"`
	RetryDelayMs int32  `json:"retryDelayMs"`
	CustomerID   string `json:"customerId"`
}

func (q *Queries) CreateWebhook(ctx context.Context, arg *CreateWebhookParams) (int64, error) {
	result, err := q.db.Exec(ctx, createWebhook,
		arg.ChainID,
		arg.Url,
		arg.MaxRetries,
		arg.TimeoutMs,
		arg.RetryDelayMs,
		arg.CustomerID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findBlockCursor = `-- name: FindBlockCursor :one
SELECT id, block_height
FROM "block_cursor"
WHERE "id" = $1
LIMIT 1
`

func (q *Queries) FindBlockCursor(ctx context.Context, chainID string) (*BlockCursor, error) {
	row := q.db.QueryRow(ctx, findBlockCursor, chainID)
	var i BlockCursor
	err := row.Scan(&i.ID, &i.BlockHeight)
	return &i, err
}

const findWebhookJobsGreaterThanID = `-- name: FindWebhookJobsGreaterThanID :many
SELECT id, created_at, chain_id, chain_url, block_height, url, max_retries, timeout_ms, retry_delay_ms 
FROM "webhook_job"
WHERE "id" > $1
ORDER BY "id" ASC
LIMIT $2
`

type FindWebhookJobsGreaterThanIDParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

func (q *Queries) FindWebhookJobsGreaterThanID(ctx context.Context, arg *FindWebhookJobsGreaterThanIDParams) ([]*WebhookJob, error) {
	rows, err := q.db.Query(ctx, findWebhookJobsGreaterThanID, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*WebhookJob{}
	for rows.Next() {
		var i WebhookJob
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.ChainID,
			&i.ChainUrl,
			&i.BlockHeight,
			&i.Url,
			&i.MaxRetries,
			&i.TimeoutMs,
			&i.RetryDelayMs,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertBlockCursor = `-- name: UpsertBlockCursor :execrows
INSERT INTO "block_cursor" (
  "id",
  "block_height"
) VALUES (
  $1,
  $2
)
ON CONFLICT ("id") DO UPDATE SET "block_height" = EXCLUDED."block_height"
`

type UpsertBlockCursorParams struct {
	ID          string `json:"id"`
	BlockHeight string `json:"blockHeight"`
}

func (q *Queries) UpsertBlockCursor(ctx context.Context, arg *UpsertBlockCursorParams) (int64, error) {
	result, err := q.db.Exec(ctx, upsertBlockCursor, arg.ID, arg.BlockHeight)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}