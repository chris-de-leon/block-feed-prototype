package blockchains

import (
	"hash/fnv"
	"strconv"
)

const (
	ETH_TESTNET_SEPOLIA ChainID = "eth-testnet-sepolia"
	ETH_TESTNET_GOERLI  ChainID = "eth-testnet-goerli"
	ETH_MAINNET         ChainID = "eth-mainnet"
	FLOW_MAINNET        ChainID = "flow-mainnet"
	FLOW_TESTNET        ChainID = "flow-testnet"
)

type (
	BlockchainOptsEnv struct {
		ConnectionURL string `env:"BLOCKCHAIN_CONNECTION_URL"`
		ChainID       string `env:"BLOCKCHAIN_ID"`
	}

	BlockchainOpts struct {
		ConnectionURL string  `validate:"required"`
		ChainID       ChainID `validate:"required"`
	}

	ChainID string
)

func GetChainID(url string) string {
	h := fnv.New32a()
	h.Write([]byte(url))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}
