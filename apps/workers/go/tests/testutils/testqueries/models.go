// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package testqueries

import (
	"time"
)

type Blockchain struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	ShardCount      int32     `json:"shardCount"`
	Url             string    `json:"url"`
	PgStoreUrl      string    `json:"pgStoreUrl"`
	RedisStoreUrl   string    `json:"redisStoreUrl"`
	RedisClusterUrl string    `json:"redisClusterUrl"`
	RedisStreamUrl  string    `json:"redisStreamUrl"`
}

type CheckoutSession struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"createdAt"`
	ClientReferenceID string    `json:"clientReferenceId"`
	SessionID         string    `json:"sessionId"`
	CustomerID        string    `json:"customerId"`
	Url               string    `json:"url"`
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
	ShardID      int32     `json:"shardId"`
}