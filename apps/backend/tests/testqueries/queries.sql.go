// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package testqueries

import (
	"context"
	"time"
)

const CreateCustomer = `-- name: CreateCustomer :execrows
INSERT INTO ` + "`" + `customer` + "`" + ` (` + "`" + `id` + "`" + `, ` + "`" + `created_at` + "`" + `) VALUES (?, DEFAULT)
`

// CreateCustomer
//
//	INSERT INTO `customer` (`id`, `created_at`) VALUES (?, DEFAULT)
func (q *Queries) CreateCustomer(ctx context.Context, id string) (int64, error) {
	result, err := q.db.ExecContext(ctx, CreateCustomer, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const CreateWebhook = `-- name: CreateWebhook :execrows
INSERT INTO ` + "`" + `webhook` + "`" + ` (` + "`" + `id` + "`" + `, ` + "`" + `created_at` + "`" + `, ` + "`" + `is_queued` + "`" + `, ` + "`" + `is_active` + "`" + `, ` + "`" + `url` + "`" + `, ` + "`" + `max_blocks` + "`" + `, ` + "`" + `max_retries` + "`" + `, ` + "`" + `timeout_ms` + "`" + `, ` + "`" + `customer_id` + "`" + `, ` + "`" + `blockchain_id` + "`" + `)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type CreateWebhookParams struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	IsQueued     bool      `json:"isQueued"`
	IsActive     bool      `json:"isActive"`
	Url          string    `json:"url"`
	MaxBlocks    int32     `json:"maxBlocks"`
	MaxRetries   int32     `json:"maxRetries"`
	TimeoutMs    int32     `json:"timeoutMs"`
	CustomerID   string    `json:"customerId"`
	BlockchainID string    `json:"blockchainId"`
}

// CreateWebhook
//
//	INSERT INTO `webhook` (`id`, `created_at`, `is_queued`, `is_active`, `url`, `max_blocks`, `max_retries`, `timeout_ms`, `customer_id`, `blockchain_id`)
//	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
func (q *Queries) CreateWebhook(ctx context.Context, arg *CreateWebhookParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, CreateWebhook,
		arg.ID,
		arg.CreatedAt,
		arg.IsQueued,
		arg.IsActive,
		arg.Url,
		arg.MaxBlocks,
		arg.MaxRetries,
		arg.TimeoutMs,
		arg.CustomerID,
		arg.BlockchainID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const CreateWebhookNodes = `-- name: CreateWebhookNodes :execrows
INSERT INTO ` + "`" + `webhook_node` + "`" + ` (` + "`" + `id` + "`" + `, ` + "`" + `created_at` + "`" + `, ` + "`" + `url` + "`" + `, ` + "`" + `blockchain_id` + "`" + `) VALUES (?, ?, ?, ?)
`

type CreateWebhookNodesParams struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Url          string    `json:"url"`
	BlockchainID string    `json:"blockchainId"`
}

// CreateWebhookNodes
//
//	INSERT INTO `webhook_node` (`id`, `created_at`, `url`, `blockchain_id`) VALUES (?, ?, ?, ?)
func (q *Queries) CreateWebhookNodes(ctx context.Context, arg *CreateWebhookNodesParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, CreateWebhookNodes,
		arg.ID,
		arg.CreatedAt,
		arg.Url,
		arg.BlockchainID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
