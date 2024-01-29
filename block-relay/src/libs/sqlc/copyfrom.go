// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: copyfrom.go

package sqlc

import (
	"context"
)

// iteratorForCreatePendingWebhookJob implements pgx.CopyFromSource.
type iteratorForCreatePendingWebhookJob struct {
	rows                 []*CreatePendingWebhookJobParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePendingWebhookJob) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePendingWebhookJob) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].BlockHeight,
		r.rows[0].ChainID,
		r.rows[0].ChainUrl,
		r.rows[0].ChannelName,
	}, nil
}

func (r iteratorForCreatePendingWebhookJob) Err() error {
	return nil
}

func (q *Queries) CreatePendingWebhookJob(ctx context.Context, arg []*CreatePendingWebhookJobParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"pending_webhook_job"}, []string{"block_height", "chain_id", "chain_url", "channel_name"}, &iteratorForCreatePendingWebhookJob{rows: arg})
}
