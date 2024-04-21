// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package sqlc

import (
	"context"
)

const ActivateWebhook = `-- name: ActivateWebhook :execrows
UPDATE ` + "`" + `webhook` + "`" + ` SET ` + "`" + `is_active` + "`" + ` = true WHERE ` + "`" + `id` + "`" + ` = ?
`

// ActivateWebhook
//
//	UPDATE `webhook` SET `is_active` = true WHERE `id` = ?
func (q *Queries) ActivateWebhook(ctx context.Context, id string) (int64, error) {
	result, err := q.db.ExecContext(ctx, ActivateWebhook, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const AssignWebhookToNode = `-- name: AssignWebhookToNode :execrows
INSERT IGNORE INTO ` + "`" + `webhook_location` + "`" + ` (` + "`" + `id` + "`" + `, ` + "`" + `created_at` + "`" + `, ` + "`" + `webhook_claim_id` + "`" + `, ` + "`" + `webhook_node_id` + "`" + `, ` + "`" + `webhook_id` + "`" + `)
VALUES (UUID(), DEFAULT, ?, ?, ?)
`

type AssignWebhookToNodeParams struct {
	WebhookClaimID string `json:"webhookClaimId"`
	WebhookNodeID  string `json:"webhookNodeId"`
	WebhookID      string `json:"webhookId"`
}

// AssignWebhookToNode
//
//	INSERT IGNORE INTO `webhook_location` (`id`, `created_at`, `webhook_claim_id`, `webhook_node_id`, `webhook_id`)
//	VALUES (UUID(), DEFAULT, ?, ?, ?)
func (q *Queries) AssignWebhookToNode(ctx context.Context, arg *AssignWebhookToNodeParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, AssignWebhookToNode, arg.WebhookClaimID, arg.WebhookNodeID, arg.WebhookID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const ClaimWebhook = `-- name: ClaimWebhook :execrows
INSERT IGNORE INTO ` + "`" + `webhook_claim` + "`" + ` (` + "`" + `id` + "`" + `, ` + "`" + `created_at` + "`" + `, ` + "`" + `claimed_by` + "`" + `, ` + "`" + `webhook_id` + "`" + `)
VALUES (UUID(), DEFAULT, ?, ?)
`

type ClaimWebhookParams struct {
	ClaimedBy string `json:"claimedBy"`
	WebhookID string `json:"webhookId"`
}

// ClaimWebhook
//
//	INSERT IGNORE INTO `webhook_claim` (`id`, `created_at`, `claimed_by`, `webhook_id`)
//	VALUES (UUID(), DEFAULT, ?, ?)
func (q *Queries) ClaimWebhook(ctx context.Context, arg *ClaimWebhookParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, ClaimWebhook, arg.ClaimedBy, arg.WebhookID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const CountWebhookNodes = `-- name: CountWebhookNodes :one
SELECT COUNT(*) FROM ` + "`" + `webhook_node` + "`" + `
`

// CountWebhookNodes
//
//	SELECT COUNT(*) FROM `webhook_node`
func (q *Queries) CountWebhookNodes(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, CountWebhookNodes)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const FindClaimedWebhook = `-- name: FindClaimedWebhook :one
SELECT webhook_claim.id, webhook_claim.created_at, webhook_claim.claimed_by, webhook_claim.webhook_id, webhook.id, webhook.created_at, webhook.is_queued, webhook.is_active, webhook.url, webhook.max_blocks, webhook.max_retries, webhook.timeout_ms, webhook.customer_id, webhook.blockchain_id
FROM ` + "`" + `webhook_claim` + "`" + ` 
LEFT JOIN ` + "`" + `webhook` + "`" + ` ON ` + "`" + `webhook` + "`" + `.` + "`" + `id` + "`" + ` = ` + "`" + `webhook_claim` + "`" + `.` + "`" + `webhook_id` + "`" + `
WHERE ` + "`" + `webhook_id` + "`" + ` = ? 
LIMIT 1
`

type FindClaimedWebhookRow struct {
	WebhookClaim WebhookClaim `json:"webhookClaim"`
	Webhook      Webhook      `json:"webhook"`
}

// FindClaimedWebhook
//
//	SELECT webhook_claim.id, webhook_claim.created_at, webhook_claim.claimed_by, webhook_claim.webhook_id, webhook.id, webhook.created_at, webhook.is_queued, webhook.is_active, webhook.url, webhook.max_blocks, webhook.max_retries, webhook.timeout_ms, webhook.customer_id, webhook.blockchain_id
//	FROM `webhook_claim`
//	LEFT JOIN `webhook` ON `webhook`.`id` = `webhook_claim`.`webhook_id`
//	WHERE `webhook_id` = ?
//	LIMIT 1
func (q *Queries) FindClaimedWebhook(ctx context.Context, webhookID string) (*FindClaimedWebhookRow, error) {
	row := q.db.QueryRowContext(ctx, FindClaimedWebhook, webhookID)
	var i FindClaimedWebhookRow
	err := row.Scan(
		&i.WebhookClaim.ID,
		&i.WebhookClaim.CreatedAt,
		&i.WebhookClaim.ClaimedBy,
		&i.WebhookClaim.WebhookID,
		&i.Webhook.ID,
		&i.Webhook.CreatedAt,
		&i.Webhook.IsQueued,
		&i.Webhook.IsActive,
		&i.Webhook.Url,
		&i.Webhook.MaxBlocks,
		&i.Webhook.MaxRetries,
		&i.Webhook.TimeoutMs,
		&i.Webhook.CustomerID,
		&i.Webhook.BlockchainID,
	)
	return &i, err
}

const GetWebhook = `-- name: GetWebhook :one
SELECT id, created_at, is_queued, is_active, url, max_blocks, max_retries, timeout_ms, customer_id, blockchain_id FROM ` + "`" + `webhook` + "`" + ` WHERE ` + "`" + `id` + "`" + ` = ?
`

// GetWebhook
//
//	SELECT id, created_at, is_queued, is_active, url, max_blocks, max_retries, timeout_ms, customer_id, blockchain_id FROM `webhook` WHERE `id` = ?
func (q *Queries) GetWebhook(ctx context.Context, id string) (*Webhook, error) {
	row := q.db.QueryRowContext(ctx, GetWebhook, id)
	var i Webhook
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.IsQueued,
		&i.IsActive,
		&i.Url,
		&i.MaxBlocks,
		&i.MaxRetries,
		&i.TimeoutMs,
		&i.CustomerID,
		&i.BlockchainID,
	)
	return &i, err
}

const LocateWebhook = `-- name: LocateWebhook :one
SELECT webhook_node.id, webhook_node.created_at, webhook_node.url, webhook_node.blockchain_id 
FROM ` + "`" + `webhook_node` + "`" + `
WHERE ` + "`" + `webhook_node` + "`" + `.` + "`" + `id` + "`" + ` IN (
  SELECT ` + "`" + `webhook_node_id` + "`" + `
  FROM ` + "`" + `webhook_location` + "`" + `
  WHERE ` + "`" + `webhook_location` + "`" + `.` + "`" + `webhook_id` + "`" + ` = ? 
)
AND ` + "`" + `webhook_node` + "`" + `.` + "`" + `blockchain_id` + "`" + ` = ?
LIMIT 1
`

type LocateWebhookParams struct {
	WebhookID    string `json:"webhookId"`
	BlockchainID string `json:"blockchainId"`
}

// LocateWebhook
//
//	SELECT webhook_node.id, webhook_node.created_at, webhook_node.url, webhook_node.blockchain_id
//	FROM `webhook_node`
//	WHERE `webhook_node`.`id` IN (
//	  SELECT `webhook_node_id`
//	  FROM `webhook_location`
//	  WHERE `webhook_location`.`webhook_id` = ?
//	)
//	AND `webhook_node`.`blockchain_id` = ?
//	LIMIT 1
func (q *Queries) LocateWebhook(ctx context.Context, arg *LocateWebhookParams) (*WebhookNode, error) {
	row := q.db.QueryRowContext(ctx, LocateWebhook, arg.WebhookID, arg.BlockchainID)
	var i WebhookNode
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Url,
		&i.BlockchainID,
	)
	return &i, err
}

const LockWebhook = `-- name: LockWebhook :one
SELECT id, created_at, is_queued, is_active, url, max_blocks, max_retries, timeout_ms, customer_id, blockchain_id FROM ` + "`" + `webhook` + "`" + ` WHERE ` + "`" + `id` + "`" + ` = ? FOR UPDATE SKIP LOCKED
`

// LockWebhook
//
//	SELECT id, created_at, is_queued, is_active, url, max_blocks, max_retries, timeout_ms, customer_id, blockchain_id FROM `webhook` WHERE `id` = ? FOR UPDATE SKIP LOCKED
func (q *Queries) LockWebhook(ctx context.Context, id string) (*Webhook, error) {
	row := q.db.QueryRowContext(ctx, LockWebhook, id)
	var i Webhook
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.IsQueued,
		&i.IsActive,
		&i.Url,
		&i.MaxBlocks,
		&i.MaxRetries,
		&i.TimeoutMs,
		&i.CustomerID,
		&i.BlockchainID,
	)
	return &i, err
}

const LockWebhookNode = `-- name: LockWebhookNode :one
SELECT webhook_node.id, webhook_node.created_at, webhook_node.url, webhook_node.blockchain_id 
FROM ` + "`" + `webhook_node` + "`" + `
LEFT JOIN (
  -- counts the number of webhooks each node is processing
  SELECT ` + "`" + `webhook_node_id` + "`" + `, COUNT(*) AS ` + "`" + `node_weight` + "`" + ` 
  FROM ` + "`" + `webhook_location` + "`" + `
  GROUP BY ` + "`" + `webhook_node_id` + "`" + `
) AS ` + "`" + `t` + "`" + ` 
ON ` + "`" + `t` + "`" + `.` + "`" + `webhook_node_id` + "`" + ` = ` + "`" + `webhook_node` + "`" + `.` + "`" + `id` + "`" + `
WHERE ` + "`" + `webhook_node` + "`" + `.` + "`" + `blockchain_id` + "`" + ` = ?
ORDER BY 
  ` + "`" + `node_weight` + "`" + ` IS NULL DESC, -- make sure nulls appear at the beginning
  ` + "`" + `node_weight` + "`" + ` ASC -- then sort in ascending order for non-null rows
LIMIT 1 
FOR UPDATE SKIP LOCKED
`

// LockWebhookNode
//
//	SELECT webhook_node.id, webhook_node.created_at, webhook_node.url, webhook_node.blockchain_id
//	FROM `webhook_node`
//	LEFT JOIN (
//	  -- counts the number of webhooks each node is processing
//	  SELECT `webhook_node_id`, COUNT(*) AS `node_weight`
//	  FROM `webhook_location`
//	  GROUP BY `webhook_node_id`
//	) AS `t`
//	ON `t`.`webhook_node_id` = `webhook_node`.`id`
//	WHERE `webhook_node`.`blockchain_id` = ?
//	ORDER BY
//	  `node_weight` IS NULL DESC, -- make sure nulls appear at the beginning
//	  `node_weight` ASC -- then sort in ascending order for non-null rows
//	LIMIT 1
//	FOR UPDATE SKIP LOCKED
func (q *Queries) LockWebhookNode(ctx context.Context, blockchainID string) (*WebhookNode, error) {
	row := q.db.QueryRowContext(ctx, LockWebhookNode, blockchainID)
	var i WebhookNode
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Url,
		&i.BlockchainID,
	)
	return &i, err
}

const UpsertBlockchain = `-- name: UpsertBlockchain :execrows
INSERT IGNORE INTO ` + "`" + `blockchain` + "`" + ` (` + "`" + `id` + "`" + `, ` + "`" + `url` + "`" + `) 
VALUES (?, ?)
ON DUPLICATE KEY UPDATE ` + "`" + `url` + "`" + ` = VALUES(` + "`" + `url` + "`" + `)
`

type UpsertBlockchainParams struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

// UpsertBlockchain
//
//	INSERT IGNORE INTO `blockchain` (`id`, `url`)
//	VALUES (?, ?)
//	ON DUPLICATE KEY UPDATE `url` = VALUES(`url`)
func (q *Queries) UpsertBlockchain(ctx context.Context, arg *UpsertBlockchainParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, UpsertBlockchain, arg.ID, arg.Url)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
