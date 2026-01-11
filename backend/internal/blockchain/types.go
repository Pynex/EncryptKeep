package blockchain

import (
	"fmt"
	"math/big"
	"time"
)

type BlockchainConfig struct {
	RPCEndpoint     string   `json:"rpc_endpoint"`
	ContractAddress string   `json:"contract_address"`
	ChainID         int64    `json:"chain_id"`
	GasLimit        uint64   `json:"gas_limit"`
	GasPrice        *big.Int `json:"gas_price"`
}

type UserData struct {
	Address     string     `json:"address"`
	Metadata    []byte     `json:"metadata"`
	EntryIDs    []*big.Int `json:"entry_ids"`
	LastUpdated time.Time  `json:"last_updated"`
}

type TransactionResult struct {
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

type SyncStatus struct {
	IsOnline     bool      `json:"is_online"`
	LastSyncTime time.Time `json:"last_sync_time"`
}

type Session struct {
	Address        string    `json:"address"`
	PrivateKey     string    `json:"private_key"`
	MasterPassword string    `json:"master_password"`
	CreatedAt      time.Time `json:"created_at"`
	LastUsed       time.Time `json:"last_used"`
}

func GetDefaultConfig() *BlockchainConfig {
	return &BlockchainConfig{
		RPCEndpoint:     "https://sepolia.base.org",
		ContractAddress: "0x02a06b3427A2D949E971Bd80606996C75ae9fEa9",
		ChainID:         84532,
		GasLimit:        1_000_000,
		GasPrice:        nil,
	}
}

func NewClientWithConfig(config *BlockchainConfig) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	return NewClient(config)
}
