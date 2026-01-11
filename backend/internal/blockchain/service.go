package blockchain

import (
	"context"
	"math/big"

	"encryptkeep-backend/internal/vault"
)

type BlockchainService interface {
	Connect() error
	Disconnect() error
	GetStatus() (*SyncStatus, error)
	IsConnected() bool

	StartSession(privateKeyHex string, masterPassword string) (*Session, error)

	StoreMetadata(ctx context.Context, data []byte) (*TransactionResult, error)
	GetUserMetadata(ctx context.Context, userAddress string) ([]byte, error)

	StoreData(ctx context.Context, data []byte) (*TransactionResult, error)
	ChangeData(ctx context.Context, dataID *big.Int, data []byte) (*TransactionResult, error)
	RemoveData(ctx context.Context, dataID *big.Int) (*TransactionResult, error)
	GetUserData(ctx context.Context, userAddress string, dataID *big.Int) ([]byte, error)
	GetActiveIds(ctx context.Context, userAddress string) ([]*big.Int, error)

	SyncVault(v *vault.LocalVault) error
}
