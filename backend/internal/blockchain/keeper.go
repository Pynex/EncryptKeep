package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type KeeperContract struct {
	client   *ethclient.Client
	contract *Keeper // Go binding
	address  common.Address
}

func NewKeeperContract(client *ethclient.Client, contractAddress string) (*KeeperContract, error) {
	address := common.HexToAddress(contractAddress)

	contract, err := NewKeeper(address, client)
	if err != nil {
		return nil, err
	}

	return &KeeperContract{
		client:   client,
		contract: contract,
		address:  address,
	}, nil
}

func (k *KeeperContract) GetUserMetadata(ctx context.Context, userAddress string) ([]byte, error) {
	data, err := k.contract.UserMetaData(&bind.CallOpts{Context: ctx}, common.HexToAddress(userAddress))
	if err != nil {
		return nil, ParseContractError(err)
	}
	return data, nil
}

func (k *KeeperContract) GetUserData(ctx context.Context, userAddress string, dataID *big.Int) ([]byte, error) {
	data, err := k.contract.UserData(&bind.CallOpts{Context: ctx}, common.HexToAddress(userAddress), dataID)
	if err != nil {
		return nil, ParseContractError(err)
	}
	return data, nil
}

func (k *KeeperContract) GetActiveIds(ctx context.Context, userAddress string) ([]*big.Int, error) {
	addr := common.HexToAddress(userAddress)
	var ids []*big.Int

	for i := int64(0); ; i++ {
		result, err := k.contract.ActiveIdsForUser(&bind.CallOpts{Context: ctx}, addr, big.NewInt(i))
		if err != nil {
			if IsContractError(err) {
				break
			}

			return nil, err
		}
		ids = append(ids, result)
	}

	return ids, nil
}

func (k *KeeperContract) StoreMetadata(ctx context.Context, auth *bind.TransactOpts, data []byte) (*TransactionResult, error) {
	tx, err := k.contract.StoreMetaData(auth, data)
	if err != nil {
		return nil, ParseContractError(err)
	}

	receipt, err := bind.WaitMined(ctx, k.client, tx)
	if err != nil {
		return nil, err
	}
	if receipt.Status == 0 {
		return nil, ParseContractError(fmt.Errorf("transaction failed"))
	}

	return &TransactionResult{
		Success:   true,
		Timestamp: time.Now(),
	}, nil
}

func (k *KeeperContract) StoreData(ctx context.Context, auth *bind.TransactOpts, data []byte) (*TransactionResult, error) {
	_, err := k.contract.StoreData(auth, data)
	if err != nil {
		return nil, ParseContractError(err)
	}
	return &TransactionResult{
		Success:   true,
		Timestamp: time.Now(),
	}, nil
}

func (k *KeeperContract) ChangeData(ctx context.Context, auth *bind.TransactOpts, dataID *big.Int, data []byte) (*TransactionResult, error) {
	tx, err := k.contract.ChangeData(auth, dataID, data)
	if err != nil {
		return nil, ParseContractError(err)
	}

	receipt, err := bind.WaitMined(ctx, k.client, tx)
	if err != nil {
		return nil, err
	}
	if receipt.Status == 0 {
		return nil, ParseContractError(fmt.Errorf("transaction failed"))
	}

	return &TransactionResult{
		Success:   true,
		Timestamp: time.Now(),
	}, nil
}

func (k *KeeperContract) RemoveData(ctx context.Context, auth *bind.TransactOpts, dataID *big.Int) (*TransactionResult, error) {
	tx, err := k.contract.RemoveData(auth, dataID)
	if err != nil {
		return nil, ParseContractError(err)
	}

	receipt, err := bind.WaitMined(ctx, k.client, tx)
	if err != nil {
		return nil, err
	}
	if receipt.Status == 0 {
		return nil, ParseContractError(fmt.Errorf("transaction failed"))
	}

	return &TransactionResult{
		Success:   true,
		Timestamp: time.Now(),
	}, nil
}
