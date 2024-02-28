// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"time"
)

type Blockchain struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

type Customer struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type Webhook struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	IsActive     bool      `json:"isActive"`
	Url          string    `json:"url"`
	MaxBlocks    int32     `json:"maxBlocks"`
	MaxRetries   int32     `json:"maxRetries"`
	TimeoutMs    int32     `json:"timeoutMs"`
	CustomerID   string    `json:"customerId"`
	BlockchainID string    `json:"blockchainId"`
}

type WebhookClaim struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ClaimedBy string    `json:"claimedBy"`
	WebhookID string    `json:"webhookId"`
}

type WebhookLocation struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	WebhookClaimID string    `json:"webhookClaimId"`
	WebhookNodeID  string    `json:"webhookNodeId"`
	WebhookID      string    `json:"webhookId"`
}

type WebhookNode struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Url          string    `json:"url"`
	BlockchainID string    `json:"blockchainId"`
}
