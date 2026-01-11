package appservice

import (
	"encryptkeep-backend/internal/blockchain"
)

type AppService struct {
	blockchainService blockchain.BlockchainService
	// vault *vault.LocalVault
}

func (as *AppService) Initialize() error {
	config := &blockchain.BlockchainConfig{
		RPCEndpoint:     "https://sepolia.infura.io/v3/",
		ContractAddress: "0x123",
		ChainID:         11155111,
		GasLimit:        10000000,
		GasPrice:        nil,
	}

	as.blockchainService = blockchain.NewBlockchainService(config)

	return as.blockchainService.Connect()
}

func (as *AppService) GetBlockchainService() blockchain.BlockchainService {
	return as.blockchainService
}
