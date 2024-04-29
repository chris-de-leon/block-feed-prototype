package e2e

import (
	"block-feed/src/libs/blockstore"
	"block-feed/src/libs/redis/redicluster"
	"block-feed/src/libs/services/processing"
	"block-feed/tests/testutils"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type (
	RequestLog struct {
		Timestamp string
		Blocks    []string
	}
)

// TODO: tests to add:
//
//	job retries / failed job handling (stream consumer)
//	fault tolerance (all)
//	add timing test (i.e. what is the average time it takes for a block to be sent to a webhook once it is sealed on the chain?)

// This test case performs the following:
//
//  1. It adds a webhook to redis cluster and schedules it for processing
//
//  2. It tests that the block streamer is correctly sending blocks to a block stream
//
//  3. It tests that a block consumer is able to (1) add blocks from the block stream to the blockstore and (2) reschedule any webhooks based on the latest block height
//
//  4. Then finally, it tests that block data is properly forwarded to a webhook URL from a webhook consumer
//
// This test is run against a live network for simplicity - in the future these test cases will be run against a local devnet
func TestBasic(t *testing.T) {
	// Defines helper constants
	const (
		FLOW_TESTNET_CHAIN_ID = "flow-testnet"
		FLOW_TESTNET_URL      = grpc.TestnetHost

		WEBHOOK_CONSUMER_NAME      = "webhook-consumer"
		WEBHOOK_CONSUMER_POOL_SIZE = 3

		BLOCK_CONSUMER_NAME       = "block-consumer"
		BLOCK_CONSUMER_BATCH_SIZE = 100

		BLOCK_FLUSH_INTERVAL_MS = 1000
		BLOCK_FLUSH_MAX_BLOCKS  = 5

		WEBHOOK_MAX_BLOCKS  = 1
		WEBHOOK_MAX_RETRIES = 3
		WEBHOOK_TIMEOUT_MS  = 5000

		TEST_DURATION_MS = 10000
		TEST_SHARDS      = 1
	)

	// Defines helper variables
	var (
		ctx    = context.Background()
		reqLog = []RequestLog{}
	)

	// Starts a mock server
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		timestamp := time.Now().UTC().String()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
			return
		} else {
			defer req.Body.Close()
		}

		blocks, err := testutils.JsonParse[[]string](string(body))
		if err != nil {
			t.Fatal(err)
		}

		reqLog = append(reqLog, RequestLog{
			Timestamp: timestamp,
			Blocks:    blocks,
		})
	}))
	t.Cleanup(func() { server.Close() })

	// Creates an error group so that we can create all containers in parallel
	containerErrGrp := new(errgroup.Group)
	var cRedisCluster *testutils.ContainerWithConnectionInfo
	var cRedisStream *testutils.ContainerWithConnectionInfo
	var cRedisStore *testutils.ContainerWithConnectionInfo
	var cTimescaleDB *testutils.ContainerWithConnectionInfo

	// Starts a timescale container
	containerErrGrp.Go(func() error {
		container, err := testutils.NewTimescaleDBContainer(ctx, t)
		if err != nil {
			return err
		} else {
			cTimescaleDB = container
		}
		return nil
	})

	// Starts a redis cluster container
	containerErrGrp.Go(func() error {
		container, err := testutils.NewRedisClusterContainer(ctx, t, testutils.REDIS_CLUSTER_MIN_NODES)
		if err != nil {
			return err
		} else {
			cRedisCluster = container
		}
		return nil
	})

	// Starts a redis container
	containerErrGrp.Go(func() error {
		container, err := testutils.NewRedisContainer(ctx, t, testutils.RedisBlockStoreCmd())
		if err != nil {
			return err
		} else {
			cRedisStore = container
		}
		return nil
	})

	// Starts a redis container
	containerErrGrp.Go(func() error {
		container, err := testutils.NewRedisContainer(ctx, t, testutils.RedisDefaultCmd())
		if err != nil {
			return err
		} else {
			cRedisStream = container
		}
		return nil
	})

	// Waits for all the containers to be created
	if err := containerErrGrp.Wait(); err != nil {
		t.Fatal(err)
	}

	// Creates a blockstore
	store, err := testutils.NewRedisOptimizedBlockStore(t, ctx,
		FLOW_TESTNET_CHAIN_ID,
		cRedisStore.Conn.Url,
		testutils.PostgresUrl(*cTimescaleDB.Conn,
			testutils.TIMESCALEDB_BLOCKSTORE_USER_UNAME,
			testutils.TIMESCALEDB_BLOCKSTORE_USER_PWORD,
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Creates a flow block streamer service
	flowBlockStreamer, err := testutils.NewFlowBlockStreamer(t, ctx,
		FLOW_TESTNET_CHAIN_ID,
		FLOW_TESTNET_URL,
		cRedisStream.Conn.Url,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Creates a block consumer service
	flowBlockConsumer, err := testutils.NewBlockStreamConsumer(t, ctx,
		TEST_SHARDS,
		store,
		FLOW_TESTNET_CHAIN_ID,
		cRedisCluster.Conn.Url,
		cRedisStream.Conn.Url,
		&processing.BlockStreamConsumerOpts{
			ConsumerName: BLOCK_CONSUMER_NAME,
			BatchSize:    BLOCK_CONSUMER_BATCH_SIZE,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Creates 1 replica of a webhook stream consumer service for each shard.
	// Each service has WEBHOOK_CONSUMER_POOL_SIZE concurrent workers.
	webhookConsumers := make([]*processing.WebhookStreamConsumer, TEST_SHARDS)
	for shardNum := range TEST_SHARDS {
		webhookConsumer, err := testutils.NewWebhookStreamConsumer(t, ctx,
			store,
			cRedisCluster.Conn.Url,
			shardNum,
			&processing.WebhookStreamConsumerOpts{
				ConsumerName: WEBHOOK_CONSUMER_NAME,
				Concurrency:  WEBHOOK_CONSUMER_POOL_SIZE,
			})
		if err != nil {
			t.Fatal(err)
		} else {
			webhookConsumers[shardNum] = webhookConsumer
		}
	}

	// Schedules a webhook for processing on each shard
	if _, err := testutils.GetTempRedisClusterClient(cRedisCluster.Conn.Url, func(client *redis.ClusterClient) (bool, error) {
		for shardNum := range TEST_SHARDS {
			if err := redicluster.NewRedisCluster(client).
				Webhooks.Set(
				ctx,
				shardNum,
				redicluster.Webhook{
					ID:           uuid.NewString(),
					URL:          server.URL,
					BlockchainID: FLOW_TESTNET_CHAIN_ID,
					MaxRetries:   WEBHOOK_MAX_RETRIES,
					MaxBlocks:    WEBHOOK_MAX_BLOCKS,
					TimeoutMs:    WEBHOOK_TIMEOUT_MS,
				},
			); err != nil {
				return false, err
			}
		}
		return true, nil
	}); err != nil {
		t.Fatal(err)
	}

	// Creates a context that will be canceled at a later time
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(TEST_DURATION_MS)*time.Millisecond)
	defer cancel()

	// Runs all services in the background
	eg := new(errgroup.Group)
	eg.Go(func() error { return flowBlockStreamer.Run(timeoutCtx) })
	eg.Go(func() error { return flowBlockConsumer.Run(timeoutCtx) })
	eg.Go(func() error {
		return store.StartFlushing(
			timeoutCtx,
			FLOW_TESTNET_CHAIN_ID,
			blockstore.RedisOptimizedBlockStoreFlushOpts{
				IntervalMs: BLOCK_FLUSH_INTERVAL_MS,
				MaxBlocks:  BLOCK_FLUSH_MAX_BLOCKS,
			},
		)
	})
	for _, consumer := range webhookConsumers {
		eg.Go(func() error { return consumer.Run(timeoutCtx) })
	}

	// Waits for the timeout (processing should occur in the background while we wait)
	// Fails the test if an unexpected error occurs
	if err := eg.Wait(); err != nil && !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "i/o timeout") {
		t.Fatal(err)
	}

	// Checks that the correct number of http calls was made
	t.Logf("Total number of requests received by webhook: %d\n", len(reqLog))
	if len(reqLog) == 0 {
		t.Fatal("Webhook received no blocks\n")
	} else {
		blockCount := 0
		for _, req := range reqLog {
			t.Logf("  => Received request at %s containing %d block(s)", req.Timestamp, len(req.Blocks))
			for _, block := range req.Blocks {
				parsedBlock, err := testutils.JsonParse[map[string]any](block)
				if err != nil {
					t.Fatal(err)
				}
				blockHeight, exists := parsedBlock["height"]
				if !exists {
					t.Fatal("block has no key named \"height\"")
				}
				blockTime, exists := parsedBlock["timestamp"]
				if !exists {
					t.Fatal("block has no key named \"timestamp\"")
				}
				t.Logf("    => Received block %.0f (block timestamp = %s)\n", blockHeight, blockTime)
			}
			blockCount += len(req.Blocks)
		}
		t.Logf("Total number of blocks received by webhook: %d\n", blockCount)
	}
}
