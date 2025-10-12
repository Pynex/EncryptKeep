package blockchain

import (
	"math/big"
	"time"
)

// BlockchainConfig содержит конфигурацию для подключения к блокчейну
type BlockchainConfig struct {
	RPCEndpoint     string   `json:"rpc_endpoint"`     // URL RPC ноды
	ContractAddress string   `json:"contract_address"` // Адрес контракта
	ChainID         int64    `json:"chain_id"`         // ID сети
	GasLimit        uint64   `json:"gas_limit"`        // Лимит газа
	GasPrice        *big.Int `json:"gas_price"`        // Цена газа
}

// UserData представляет данные пользователя в блокчейне
type UserData struct {
	Address     string     `json:"address"`      // Адрес пользователя
	Metadata    []byte     `json:"metadata"`     // Метаданные пользователя (один blob)
	EntryIDs    []*big.Int `json:"entry_ids"`    // ID записей из activeIdsForUser
	LastUpdated time.Time  `json:"last_updated"` // Время последнего обновления
}

// TransactionResult результат транзакции
type TransactionResult struct {
	Success   bool      `json:"success"`   // Успешность транзакции
	Timestamp time.Time `json:"timestamp"` // Время транзакции
}

// SyncStatus статус синхронизации
type SyncStatus struct {
	IsOnline     bool      `json:"is_online"`      // Онлайн ли блокчейн
	LastSyncTime time.Time `json:"last_sync_time"` // Время последней синхронизации
}

// Session представляет сессию пользователя
type Session struct {
	Address    string    `json:"address"`     // Адрес пользователя (из приватного ключа)
	PrivateKey string    `json:"private_key"` // Приватный ключ (хранится в памяти)
	MasterKey  []byte    `json:"master_key"`  // Мастер-ключ (деривированный из пароля)
	CreatedAt  time.Time `json:"created_at"`  // Время создания сессии
	LastUsed   time.Time `json:"last_used"`   // Время последнего использования
}
