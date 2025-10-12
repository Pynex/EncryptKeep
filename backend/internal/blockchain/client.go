package blockchain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client представляет клиент для работы с блокчейном
type Client struct {
	config   *BlockchainConfig
	client   *ethclient.Client
	contract *KeeperContract
	chainID  *big.Int
	session  *Session // Текущая сессия пользователя
}

// NewClient создаёт новый клиент блокчейна
func NewClient(config *BlockchainConfig) (*Client, error) {
	// Подключение к RPC ноде
	client, err := ethclient.Dial(config.RPCEndpoint)
	if err != nil {
		return nil, ErrConnectionFailed
	}

	// Получение chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, ErrConnectionFailed
	}

	// Проверка chainID
	if chainID.Int64() != config.ChainID {
		return nil, ErrInvalidChainID
	}

	// Создание контракта
	contract, err := NewKeeperContract(client, config.ContractAddress)
	if err != nil {
		return nil, ErrContractNotFound
	}

	return &Client{
		config:   config,
		client:   client,
		contract: contract,
		chainID:  chainID,
	}, nil
}

// GetUserMetadata получает метаданные пользователя
func (c *Client) GetUserMetadata(ctx context.Context, userAddress string) ([]byte, error) {
	return c.contract.GetUserMetadata(ctx, userAddress)
}

// GetUserData получает данные пользователя по ID
func (c *Client) GetUserData(ctx context.Context, userAddress string, dataID *big.Int) ([]byte, error) {
	return c.contract.GetUserData(ctx, userAddress, dataID)
}

// GetActiveIds получает список активных ID для пользователя
func (c *Client) GetActiveIds(ctx context.Context, userAddress string) ([]*big.Int, error) {
	return c.contract.GetActiveIds(ctx, userAddress)
}

// StoreMetadata сохраняет метаданные пользователя
func (c *Client) StoreMetadata(ctx context.Context, data []byte) (*TransactionResult, error) {
	if c.session == nil {
		return nil, ErrInvalidPrivateKey
	}

	auth, err := c.createAuth(c.session.PrivateKey)
	if err != nil {
		return nil, err
	}

	return c.contract.StoreMetadata(ctx, auth, data)
}

// StoreData сохраняет данные пользователя
func (c *Client) StoreData(ctx context.Context, data []byte) (*TransactionResult, error) {
	if c.session == nil {
		return nil, ErrInvalidPrivateKey
	}

	auth, err := c.createAuth(c.session.PrivateKey)
	if err != nil {
		return nil, err
	}

	return c.contract.StoreData(ctx, auth, data)
}

// ChangeData изменяет данные пользователя
func (c *Client) ChangeData(ctx context.Context, dataID *big.Int, data []byte) (*TransactionResult, error) {
	if c.session == nil {
		return nil, ErrInvalidPrivateKey
	}

	auth, err := c.createAuth(c.session.PrivateKey)
	if err != nil {
		return nil, err
	}

	return c.contract.ChangeData(ctx, auth, dataID, data)
}

// RemoveData удаляет данные пользователя
func (c *Client) RemoveData(ctx context.Context, dataID *big.Int) (*TransactionResult, error) {
	if c.session == nil {
		return nil, ErrInvalidPrivateKey
	}

	auth, err := c.createAuth(c.session.PrivateKey)
	if err != nil {
		return nil, err
	}

	return c.contract.RemoveData(ctx, auth, dataID)
}

// createAuth создаёт авторизацию для транзакций
func (c *Client) createAuth(privateKeyHex string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, ErrInvalidPrivateKey
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrInvalidPrivateKey
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, c.chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = c.config.GasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

// Close закрывает соединение с блокчейном
func (c *Client) Close() error {
	if c.client != nil {
		c.client.Close()
	}
	return nil
}

// GetSyncStatus получает статус синхронизации
func (c *Client) GetSyncStatus(ctx context.Context) (*SyncStatus, error) {
	// Проверка подключения
	_, err := c.client.NetworkID(ctx)
	isOnline := err == nil

	return &SyncStatus{
		IsOnline:     isOnline,
		LastSyncTime: time.Now(),
	}, nil
}

// CreateSession создаёт новую сессию пользователя
func (c *Client) CreateSession(privateKeyHex string, masterKey []byte) (*Session, error) {
	// Получение адреса из приватного ключа
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, ErrInvalidPrivateKey
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrInvalidPrivateKey
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	session := &Session{
		Address:    address.Hex(),
		PrivateKey: privateKeyHex,
		MasterKey:  masterKey,
		CreatedAt:  time.Now(),
		LastUsed:   time.Now(),
	}

	c.session = session
	return session, nil
}

// GetSession получает текущую сессию
func (c *Client) GetSession() *Session {
	if c.session != nil {
		c.session.LastUsed = time.Now()
	}
	return c.session
}

// ClearSession очищает текущую сессию
func (c *Client) ClearSession() {
	if c.session != nil {
		// Очищаем чувствительные данные
		c.session.PrivateKey = ""
		c.session.MasterKey = nil
		c.session = nil
	}
}
