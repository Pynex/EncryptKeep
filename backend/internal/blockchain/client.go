package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	config   *BlockchainConfig
	client   *ethclient.Client
	contract *KeeperContract
	chainID  *big.Int
	session  *Session
}

// type GasConfig struct {
// 	MaxGasPrice   *big.Int
// 	GasMultiplier float64
// 	MaxGasLimit   uint64
// 	PriorityFee   *big.Int
// 	MaxFeePerGas  *big.Int
// }

// type ClientConfig struct {
// 	RPCEndpoint string
// 	ContractAddress string
// 	ChainID int64
// }

func NewClient(config *BlockchainConfig) (*Client, error) {
	client, err := ethclient.Dial(config.RPCEndpoint)
	if err != nil {
		return nil, ErrConnectionFailed
	}

	contract, err := NewKeeperContract(client, config.ContractAddress)
	if err != nil {
		return nil, ErrContractNotFound
	}

	return &Client{
		config:   config,
		client:   client,
		contract: contract,
		chainID:  big.NewInt(config.ChainID),
	}, nil
}

func (c *Client) EstimateGas(ctx context.Context, method string, args ...interface{}) (uint64, error) {
	if method == "storeMetadata" {
		data := args[0].([]byte)
		tx, err := c.contract.contract.StoreMetaData(&bind.TransactOpts{
			Context: ctx,
			NoSend:  true,
		}, data)
		if err != nil {
			return 0, err
		}

		gas, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
			To:   &c.contract.address,
			Data: tx.Data(),
		})
		if err != nil {
			return 0, err
		}

		return gas, nil
	}

	if method == "storeData" {
		data := args[0].([]byte)
		tx, err := c.contract.contract.StoreData(&bind.TransactOpts{
			Context: ctx,
			NoSend:  true,
		}, data)
		if err != nil {
			return 0, err
		}

		gas, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
			To:   &c.contract.address,
			Data: tx.Data(),
		})
		if err != nil {
			return 0, err
		}

		return gas, nil
	}

	if method == "changeData" {
		dataID := args[0].(*big.Int)
		data := args[1].([]byte)
		tx, err := c.contract.contract.ChangeData(&bind.TransactOpts{
			Context: ctx,
			NoSend:  true,
		}, dataID, data)
		if err != nil {
			return 0, err
		}

		gas, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
			To:   &c.contract.address,
			Data: tx.Data(),
		})
		if err != nil {
			return 0, err
		}

		return gas, nil
	}

	if method == "removeData" {
		dataID := args[0].(*big.Int)
		tx, err := c.contract.contract.RemoveData(&bind.TransactOpts{
			Context: ctx,
			NoSend:  true,
		}, dataID)
		if err != nil {
			return 0, err
		}

		gas, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
			To:   &c.contract.address,
			Data: tx.Data(),
		})
		if err != nil {
			return 0, err
		}

		return gas, nil
	}

	return 0, fmt.Errorf("unknown method")
}

// func (c *Client) GetOptimalGasPrice(ctx context.Context) (*big.Int, *big.Int, error) {
// 	return c.contract
// }

func (c *Client) GetUserMetadata(ctx context.Context, userAddress string) ([]byte, error) {
	return c.contract.GetUserMetadata(ctx, userAddress)
}

func (c *Client) GetUserData(ctx context.Context, userAddress string, dataID *big.Int) ([]byte, error) {
	return c.contract.GetUserData(ctx, userAddress, dataID)
}

func (c *Client) GetActiveIds(ctx context.Context, userAddress string) ([]*big.Int, error) {
	return c.contract.GetActiveIds(ctx, userAddress)
}

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

func (c *Client) Close() error {
	if c.client != nil {
		c.client.Close()
	}
	return nil
}

func (c *Client) GetSyncStatus(ctx context.Context) (*SyncStatus, error) {
	_, err := c.client.NetworkID(ctx)
	isOnline := err == nil

	return &SyncStatus{
		IsOnline:     isOnline,
		LastSyncTime: time.Now(),
	}, nil
}

func (c *Client) CreateSession(privateKeyHex string, masterPassword string) (*Session, error) {
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
		Address:        address.Hex(),
		PrivateKey:     privateKeyHex,
		MasterPassword: masterPassword,
		CreatedAt:      time.Now(),
		LastUsed:       time.Now(),
	}

	c.session = session
	return session, nil
}

func (c *Client) GetSession() *Session {
	if c.session != nil {
		c.session.LastUsed = time.Now()
	}
	return c.session
}

func (c *Client) ClearSession() {
	if c.session != nil {
		c.session.PrivateKey = ""
		c.session.MasterPassword = ""
		c.session = nil
	}
}
