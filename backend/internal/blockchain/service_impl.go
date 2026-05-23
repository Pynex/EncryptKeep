package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"encryptkeep-backend/internal/codec"
	"encryptkeep-backend/internal/vault"
)

type BlockchainServiceImpl struct {
	client *Client
	config *BlockchainConfig
}

func NewBlockchainService(config *BlockchainConfig) *BlockchainServiceImpl {
	return &BlockchainServiceImpl{
		config: config,
	}
}

func (bs *BlockchainServiceImpl) Connect() error {
	client, err := NewClient(bs.config)
	if err != nil {
		return err
	}

	bs.client = client

	return nil
}

func (bs *BlockchainServiceImpl) GetStatus() (*SyncStatus, error) {
	if bs.client == nil {
		return &SyncStatus{IsOnline: false}, nil
	}

	return bs.client.GetSyncStatus(context.Background())
}

func (bs *BlockchainServiceImpl) Disconnect() error {
	if bs.client != nil {
		return bs.client.Close()
	}

	return nil
}

func (bs *BlockchainServiceImpl) IsConnected() bool {
	return bs.client != nil
}

func (bs *BlockchainServiceImpl) StartSession(privateKeyHex string, masterPassword string) (*Session, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.CreateSession(privateKeyHex, masterPassword)
}

func (bs *BlockchainServiceImpl) StoreMetadata(ctx context.Context, data []byte) (*TransactionResult, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.StoreMetadata(ctx, data)
}

func (bs *BlockchainServiceImpl) GetUserMetadata(ctx context.Context, userAddress string) ([]byte, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.GetUserMetadata(ctx, userAddress)
}

func (bs *BlockchainServiceImpl) StoreData(ctx context.Context, data []byte) (*TransactionResult, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.StoreData(ctx, data)
}

func (bs *BlockchainServiceImpl) ChangeData(ctx context.Context, dataID *big.Int, data []byte) (*TransactionResult, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.ChangeData(ctx, dataID, data)
}

func (bs *BlockchainServiceImpl) RemoveData(ctx context.Context, dataID *big.Int) (*TransactionResult, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.RemoveData(ctx, dataID)
}

func (bs *BlockchainServiceImpl) GetUserData(ctx context.Context, userAddress string, dataID *big.Int) ([]byte, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.GetUserData(ctx, userAddress, dataID)
}

func (bs *BlockchainServiceImpl) GetActiveIds(ctx context.Context, userAddress string) ([]*big.Int, error) {
	if bs.client == nil {
		return nil, ErrNotConnected
	}
	return bs.client.GetActiveIds(ctx, userAddress)
}

func (bs *BlockchainServiceImpl) SyncVault(v *vault.LocalVault) error {
	if bs.client == nil {
		return ErrNotConnected
	}

	session := bs.client.GetSession()
	if session == nil || session.MasterPassword == "" {
		return ErrInvalidPrivateKey
	}
	ctx := context.Background()
	userAddr := session.Address

	var metaBytes []byte
	err := bs.callWithReconnect(func() error {
		b, callErr := bs.client.GetUserMetadata(ctx, userAddr)
		if callErr != nil {
			return callErr
		}
		metaBytes = b
		return nil
	})
	if err != nil {
		return err
	}

	cdc := codec.NewCodec()
	meta := &vault.UserMetadata{
		Version:      "1.0",
		Settings:     map[string]string{},
		PasswordIDs:  []string{},
		UpdatedAt:    time.Now(),
		TotalEntries: 0,
	}
	if len(metaBytes) > 0 {
		if decoded, err := cdc.UnpackMetadata(metaBytes, session.MasterPassword); err == nil {
			meta = decoded
		} else {
			return fmt.Errorf("decrypt vault metadata: %w", err)
		}
	}

	var ids []*big.Int
	err = bs.callWithReconnect(func() error {
		fetchedIDs, callErr := bs.client.GetActiveIds(ctx, userAddr)
		if callErr != nil {
			return callErr
		}
		ids = fetchedIDs
		return nil
	})
	if err != nil {
		return err
	}

	entries := make(map[string]*vault.PasswordEntry)
	blockchainEntries := make(map[string]*big.Int)

	for _, id := range ids {
		if id == nil {
			continue
		}
		var dataBytes []byte
		err := bs.callWithReconnect(func() error {
			fetched, callErr := bs.client.GetUserData(ctx, userAddr, id)
			if callErr != nil {
				return callErr
			}
			dataBytes = fetched
			return nil
		})
		if err != nil {
			return err
		}

		entry, err := cdc.UnpackEntry(dataBytes, session.MasterPassword)
		if err != nil {
			return fmt.Errorf("decrypt vault entry (contract id %s): %w", id.String(), err)
		}

		entries[entry.ID] = entry
		blockchainEntries[entry.ID] = id
	}

	v.Metadata = meta
	v.Entries = entries
	v.BlockchainEntries = blockchainEntries
	v.LastSyncTime = meta.UpdatedAt
	v.IsDirty = false

	return nil
}

func (bs *BlockchainServiceImpl) callWithReconnect(call func() error) error {
	err := call()
	if err == nil {
		return nil
	}
	if !isTransientRPCError(err) {
		return err
	}
	if reconnectErr := bs.reconnectAndRestoreSession(); reconnectErr != nil {
		return fmt.Errorf("rpc call failed: %w (reconnect failed: %v)", err, reconnectErr)
	}
	return call()
}

func (bs *BlockchainServiceImpl) reconnectAndRestoreSession() error {
	if bs.client == nil {
		return ErrNotConnected
	}
	session := bs.client.GetSession()
	if session == nil || session.PrivateKey == "" {
		return ErrInvalidPrivateKey
	}
	_ = bs.client.Close()

	client, err := NewClient(bs.config)
	if err != nil {
		return err
	}
	bs.client = client
	_, err = bs.client.CreateSession(session.PrivateKey, session.MasterPassword)
	return err
}

func isTransientRPCError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	markers := []string{
		"timeout",
		"deadline exceeded",
		"temporarily unavailable",
		"connection reset",
		"connection refused",
		"connection aborted",
		"no such host",
		"getaddrinfo",
		"wsarecv",
		"eof",
	}
	for _, marker := range markers {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}
