// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

type CacheBlocksParams struct {
	BlockchainID string `json:"blockchainId"`
	BlockHeight  int64  `json:"blockHeight"`
	Block        []byte `json:"block"`
}

const CreateWebhook = `-- name: CreateWebhook :execrows

WITH 
  inserted_user AS (
    INSERT INTO "customer" ("id", "created_at")
    VALUES ($2::TEXT, DEFAULT)
    ON CONFLICT ("id") DO NOTHING
  ), 
  inserted_blockchain AS (
    INSERT INTO "blockchain" ("id", "url") 
    VALUES ($3, $4)
    ON CONFLICT ("id") DO UPDATE SET "url" = EXCLUDED."url"
  ),
  inserted_webhook AS (
    INSERT INTO "webhook" (
      "id",
      "created_at",
      "url",
      "max_blocks",
      "max_retries",
      "timeout_ms",
      "customer_id",
      "blockchain_id"
    ) VALUES (
      DEFAULT,
      DEFAULT,
      $5,
      $6,
      $7,
      $8,
      $2,
      $3
    )
    RETURNING "id"
  )
INSERT INTO "webhook_job" ("id", "created_at", "block_height", "webhook_id") 
VALUES (DEFAULT, DEFAULT, $1, (SELECT "id" FROM inserted_webhook))
`

type CreateWebhookParams struct {
	LatestBlockHeight int64  `json:"latestBlockHeight"`
	CustomerID        string `json:"customerId"`
	BlockchainID      string `json:"blockchainId"`
	BlockchainUrl     string `json:"blockchainUrl"`
	Url               string `json:"url"`
	MaxBlocks         int32  `json:"maxBlocks"`
	MaxRetries        int32  `json:"maxRetries"`
	TimeoutMs         int32  `json:"timeoutMs"`
}

// TODO: remove this since it is only used for testing purposes
//
//	WITH
//	  inserted_user AS (
//	    INSERT INTO "customer" ("id", "created_at")
//	    VALUES ($2::TEXT, DEFAULT)
//	    ON CONFLICT ("id") DO NOTHING
//	  ),
//	  inserted_blockchain AS (
//	    INSERT INTO "blockchain" ("id", "url")
//	    VALUES ($3, $4)
//	    ON CONFLICT ("id") DO UPDATE SET "url" = EXCLUDED."url"
//	  ),
//	  inserted_webhook AS (
//	    INSERT INTO "webhook" (
//	      "id",
//	      "created_at",
//	      "url",
//	      "max_blocks",
//	      "max_retries",
//	      "timeout_ms",
//	      "customer_id",
//	      "blockchain_id"
//	    ) VALUES (
//	      DEFAULT,
//	      DEFAULT,
//	      $5,
//	      $6,
//	      $7,
//	      $8,
//	      $2,
//	      $3
//	    )
//	    RETURNING "id"
//	  )
//	INSERT INTO "webhook_job" ("id", "created_at", "block_height", "webhook_id")
//	VALUES (DEFAULT, DEFAULT, $1, (SELECT "id" FROM inserted_webhook))
func (q *Queries) CreateWebhook(ctx context.Context, arg *CreateWebhookParams) (int64, error) {
	result, err := q.db.Exec(ctx, CreateWebhook,
		arg.LatestBlockHeight,
		arg.CustomerID,
		arg.BlockchainID,
		arg.BlockchainUrl,
		arg.Url,
		arg.MaxBlocks,
		arg.MaxRetries,
		arg.TimeoutMs,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const GetLatestCachedBlockHeight = `-- name: GetLatestCachedBlockHeight :one
SELECT "blockchain_id", "block_height" 
FROM "block_cache" 
WHERE "blockchain_id" = $1
ORDER BY "block_height" DESC
LIMIT 1
`

type GetLatestCachedBlockHeightRow struct {
	BlockchainID string `json:"blockchainId"`
	BlockHeight  int64  `json:"blockHeight"`
}

// GetLatestCachedBlockHeight
//
//	SELECT "blockchain_id", "block_height"
//	FROM "block_cache"
//	WHERE "blockchain_id" = $1
//	ORDER BY "block_height" DESC
//	LIMIT 1
func (q *Queries) GetLatestCachedBlockHeight(ctx context.Context, blockchainID string) (*GetLatestCachedBlockHeightRow, error) {
	row := q.db.QueryRow(ctx, GetLatestCachedBlockHeight, blockchainID)
	var i GetLatestCachedBlockHeightRow
	err := row.Scan(&i.BlockchainID, &i.BlockHeight)
	return &i, err
}

const GetWebhookJob = `-- name: GetWebhookJob :one
SELECT 
  "webhook_job"."id" AS "id",
  "webhook"."id" AS "webhook_id",
  "webhook"."url" AS "webhook_url",
  "webhook"."timeout_ms" AS "webhook_timeout_ms",
  (
    SELECT array_agg("block")::JSONB[] 
    FROM "block_cache" 
    WHERE "block_cache"."blockchain_id" = "webhook"."blockchain_id"
    AND "block_cache"."block_height" >= "webhook_job"."block_height"
    LIMIT "webhook"."max_blocks"
  ) AS "cached_blocks"
FROM "webhook_job"
INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
WHERE "webhook_job"."id" = $1
LIMIT 1
`

type GetWebhookJobRow struct {
	ID               int64     `json:"id"`
	WebhookID        uuid.UUID `json:"webhookId"`
	WebhookUrl       string    `json:"webhookUrl"`
	WebhookTimeoutMs int32     `json:"webhookTimeoutMs"`
	CachedBlocks     [][]byte  `json:"cachedBlocks"`
}

// GetWebhookJob
//
//	SELECT
//	  "webhook_job"."id" AS "id",
//	  "webhook"."id" AS "webhook_id",
//	  "webhook"."url" AS "webhook_url",
//	  "webhook"."timeout_ms" AS "webhook_timeout_ms",
//	  (
//	    SELECT array_agg("block")::JSONB[]
//	    FROM "block_cache"
//	    WHERE "block_cache"."blockchain_id" = "webhook"."blockchain_id"
//	    AND "block_cache"."block_height" >= "webhook_job"."block_height"
//	    LIMIT "webhook"."max_blocks"
//	  ) AS "cached_blocks"
//	FROM "webhook_job"
//	INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
//	WHERE "webhook_job"."id" = $1
//	LIMIT 1
func (q *Queries) GetWebhookJob(ctx context.Context, id int64) (*GetWebhookJobRow, error) {
	row := q.db.QueryRow(ctx, GetWebhookJob, id)
	var i GetWebhookJobRow
	err := row.Scan(
		&i.ID,
		&i.WebhookID,
		&i.WebhookUrl,
		&i.WebhookTimeoutMs,
		&i.CachedBlocks,
	)
	return &i, err
}

const GetWebhookJobByWebhookID = `-- name: GetWebhookJobByWebhookID :one
SELECT id, created_at, block_height, webhook_id FROM "webhook_job" WHERE "webhook_id" = $1
`

// GetWebhookJobByWebhookID
//
//	SELECT id, created_at, block_height, webhook_id FROM "webhook_job" WHERE "webhook_id" = $1
func (q *Queries) GetWebhookJobByWebhookID(ctx context.Context, webhookID uuid.UUID) (*WebhookJob, error) {
	row := q.db.QueryRow(ctx, GetWebhookJobByWebhookID, webhookID)
	var i WebhookJob
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.BlockHeight,
		&i.WebhookID,
	)
	return &i, err
}

const GetWebhookJobs = `-- name: GetWebhookJobs :many
SELECT 
  "webhook_job"."id" AS "id",
  "webhook"."id" AS "webhook_id",
  "webhook"."max_retries" AS "webhook_max_retries"
FROM "webhook_job"
INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
WHERE "webhook_job"."id" > $1
ORDER BY "webhook_job"."id" ASC
LIMIT $2
`

type GetWebhookJobsParams struct {
	CursorID int64 `json:"cursorId"`
	Limit    int32 `json:"limit"`
}

type GetWebhookJobsRow struct {
	ID                int64     `json:"id"`
	WebhookID         uuid.UUID `json:"webhookId"`
	WebhookMaxRetries int32     `json:"webhookMaxRetries"`
}

// GetWebhookJobs
//
//	SELECT
//	  "webhook_job"."id" AS "id",
//	  "webhook"."id" AS "webhook_id",
//	  "webhook"."max_retries" AS "webhook_max_retries"
//	FROM "webhook_job"
//	INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
//	WHERE "webhook_job"."id" > $1
//	ORDER BY "webhook_job"."id" ASC
//	LIMIT $2
func (q *Queries) GetWebhookJobs(ctx context.Context, arg *GetWebhookJobsParams) ([]*GetWebhookJobsRow, error) {
	rows, err := q.db.Query(ctx, GetWebhookJobs, arg.CursorID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetWebhookJobsRow{}
	for rows.Next() {
		var i GetWebhookJobsRow
		if err := rows.Scan(&i.ID, &i.WebhookID, &i.WebhookMaxRetries); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const RescheduleWebhookJob = `-- name: RescheduleWebhookJob :execrows
WITH webhook_jobs_with_same_blockchain AS (
  -- Gets all webhook jobs that have a blockchain ID that's identical 
  -- to the job being rescheduled.
  SELECT "webhook_job"."id", "webhook_job"."block_height", "webhook"."blockchain_id"
  FROM "webhook_job"
  INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
  WHERE "webhook"."blockchain_id" IN (
    SELECT "webhook"."blockchain_id" AS "id"
    FROM "webhook_job"
    INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
    WHERE "webhook_job"."id" = $2
    LIMIT 1
  )
),
clean_block_cache AS (
  -- To clean the cache we first need to filter out the blocks that 
  -- are associated with the same chain as job being rescheduled. Once
  -- we know which blocks these are, we find the job with the smallest 
  -- block height, and any cached block that has a height smaller than 
  -- this is no longer going to be queried and can safely be deleted.  
  DELETE FROM "block_cache"
  WHERE "block_cache"."blockchain_id" IN (
    SELECT DISTINCT "blockchain_id" 
    FROM webhook_jobs_with_same_blockchain
  )
  AND "block_cache"."block_height" < (
    SELECT MIN("block_height")
    FROM webhook_jobs_with_same_blockchain
  )
),
deleted_job AS (
  -- Deletes the webhook job so that we can generate a new auto-
  -- incrementing ID for it
  DELETE FROM "webhook_job" WHERE "id" = $2 RETURNING "webhook_id" 
)
INSERT INTO "webhook_job" ("id", "created_at", "block_height", "webhook_id") 
VALUES (DEFAULT, DEFAULT, $1, (SELECT "webhook_id" FROM deleted_job))
ON CONFLICT ("webhook_id") DO NOTHING
`

type RescheduleWebhookJobParams struct {
	BlockHeight int64 `json:"blockHeight"`
	ID          int64 `json:"id"`
}

// Creates a new job with an updated block height
//
//	WITH webhook_jobs_with_same_blockchain AS (
//	  -- Gets all webhook jobs that have a blockchain ID that's identical
//	  -- to the job being rescheduled.
//	  SELECT "webhook_job"."id", "webhook_job"."block_height", "webhook"."blockchain_id"
//	  FROM "webhook_job"
//	  INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
//	  WHERE "webhook"."blockchain_id" IN (
//	    SELECT "webhook"."blockchain_id" AS "id"
//	    FROM "webhook_job"
//	    INNER JOIN "webhook" ON "webhook"."id" = "webhook_job"."webhook_id"
//	    WHERE "webhook_job"."id" = $2
//	    LIMIT 1
//	  )
//	),
//	clean_block_cache AS (
//	  -- To clean the cache we first need to filter out the blocks that
//	  -- are associated with the same chain as job being rescheduled. Once
//	  -- we know which blocks these are, we find the job with the smallest
//	  -- block height, and any cached block that has a height smaller than
//	  -- this is no longer going to be queried and can safely be deleted.
//	  DELETE FROM "block_cache"
//	  WHERE "block_cache"."blockchain_id" IN (
//	    SELECT DISTINCT "blockchain_id"
//	    FROM webhook_jobs_with_same_blockchain
//	  )
//	  AND "block_cache"."block_height" < (
//	    SELECT MIN("block_height")
//	    FROM webhook_jobs_with_same_blockchain
//	  )
//	),
//	deleted_job AS (
//	  -- Deletes the webhook job so that we can generate a new auto-
//	  -- incrementing ID for it
//	  DELETE FROM "webhook_job" WHERE "id" = $2 RETURNING "webhook_id"
//	)
//	INSERT INTO "webhook_job" ("id", "created_at", "block_height", "webhook_id")
//	VALUES (DEFAULT, DEFAULT, $1, (SELECT "webhook_id" FROM deleted_job))
//	ON CONFLICT ("webhook_id") DO NOTHING
func (q *Queries) RescheduleWebhookJob(ctx context.Context, arg *RescheduleWebhookJobParams) (int64, error) {
	result, err := q.db.Exec(ctx, RescheduleWebhookJob, arg.BlockHeight, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const UpsertBlockchain = `-- name: UpsertBlockchain :execrows
INSERT INTO "blockchain" ("id", "url") 
VALUES ($1, $2)
ON CONFLICT ("id") DO UPDATE SET "url" = EXCLUDED."url"
`

type UpsertBlockchainParams struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

// UpsertBlockchain
//
//	INSERT INTO "blockchain" ("id", "url")
//	VALUES ($1, $2)
//	ON CONFLICT ("id") DO UPDATE SET "url" = EXCLUDED."url"
func (q *Queries) UpsertBlockchain(ctx context.Context, arg *UpsertBlockchainParams) (int64, error) {
	result, err := q.db.Exec(ctx, UpsertBlockchain, arg.ID, arg.Url)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
