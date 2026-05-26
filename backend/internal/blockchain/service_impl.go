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

const rpcCallTimeout = 12 * time.Second

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
	userAddr := session.Address

	var metaBytes []byte
	err := bs.callWithReconnect(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), rpcCallTimeout)
		defer cancel()
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
		ctx, cancel := context.WithTimeout(context.Background(), rpcCallTimeout)
		defer cancel()
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
			ctx, cancel := context.WithTimeout(context.Background(), rpcCallTimeout)
			defer cancel()
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
	failedEndpoints := make([]string, 0, len(RPCEndpoints(bs.config)))
	maxAttempts := len(RPCEndpoints(bs.config)) + 1
	if maxAttempts < 2 {
		maxAttempts = 2
	}

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err := call()
		if err == nil {
			return nil
		}
		lastErr = err
		if !isTransientRPCError(err) {
			return wrapRPCError(bs.client, err)
		}

		if bs.client != nil {
			failedEndpoints = appendUniqueEndpoint(failedEndpoints, bs.client.ActiveRPCEndpoint())
		}
		if reconnectErr := bs.reconnectExcluding(failedEndpoints...); reconnectErr != nil {
			lastErr = fmt.Errorf("%w (reconnect failed: %v)", err, reconnectErr)
			continue
		}
	}
	return wrapRPCError(bs.client, fmt.Errorf("all rpc endpoints exhausted: %w", lastErr))
}

func (bs *BlockchainServiceImpl) reconnectExcluding(exclude ...string) error {
	if bs.client == nil {
		return ErrNotConnected
	}
	session := bs.client.GetSession()
	if session == nil || session.PrivateKey == "" {
		return ErrInvalidPrivateKey
	}
	_ = bs.client.Close()

	client, err := NewClientExcluding(bs.config, exclude...)
	if err != nil {
		return err
	}
	bs.client = client
	_, err = bs.client.CreateSession(session.PrivateKey, session.MasterPassword)
	return err
}

func appendUniqueEndpoint(list []string, endpoint string) []string {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return list
	}
	for _, existing := range list {
		if existing == endpoint {
			return list
		}
	}
	return append(list, endpoint)
}

func wrapRPCError(client *Client, err error) error {
	if err == nil {
		return nil
	}
	if client == nil {
		return err
	}
	if ep := client.ActiveRPCEndpoint(); ep != "" {
		return fmt.Errorf("rpc %s: %w", ep, err)
	}
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
		"failed to respond",
		"i/o timeout",
		"network is unreachable",
		"eof",
	}
	for _, marker := range markers {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}
