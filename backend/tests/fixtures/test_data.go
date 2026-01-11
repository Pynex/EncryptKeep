package fixtures

import (
	"math/big"
	"time"

	"encryptkeep-backend/internal/blockchain"
	"encryptkeep-backend/internal/vault"
)

// TestPrivateKey - тестовый приватный ключ (НЕ ИСПОЛЬЗУЙ В ПРОДАКШЕНЕ!)
const TestPrivateKey = "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

// TestMasterPassword - тестовый мастер-пароль
const TestMasterPassword = "testpassword123"

// TestUserAddress - тестовый адрес пользователя
const TestUserAddress = "0x1234567890123456789012345678901234567890"

// TestContractAddress - тестовый адрес контракта
const TestContractAddress = "0xabcdef1234567890abcdef1234567890abcdef12"

// TestRPCEndpoint - тестовый RPC endpoint
const TestRPCEndpoint = "http://localhost:8545"

// TestChainID - тестовый Chain ID
const TestChainID = int64(1337)

// GetTestPasswordEntry создает тестовую запись пароля
func GetTestPasswordEntry() *vault.PasswordEntry {
	return &vault.PasswordEntry{
		ID:         "test-entry-1",
		Title:      "Test Website",
		Username:   "testuser",
		Password:   "testpassword123",
		URL:        "https://example.com",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsFavorite: false,
	}
}

// GetTestPasswordEntries создает несколько тестовых записей паролей
func GetTestPasswordEntries() []*vault.PasswordEntry {
	return []*vault.PasswordEntry{
		{
			ID:         "test-entry-1",
			Title:      "Test Website 1",
			Username:   "user1",
			Password:   "password1",
			URL:        "https://example1.com",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			IsFavorite: false,
		},
		{
			ID:         "test-entry-2",
			Title:      "Test Website 2",
			Username:   "user2",
			Password:   "password2",
			URL:        "https://example2.com",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			IsFavorite: true,
		},
		{
			ID:         "test-entry-3",
			Title:      "Test Website 3",
			Username:   "user3",
			Password:   "password3",
			URL:        "https://example3.com",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			IsFavorite: false,
		},
	}
}

// GetTestUserMetadata создает тестовые метаданные пользователя
func GetTestUserMetadata() *vault.UserMetadata {
	return &vault.UserMetadata{
		Version:      "1.0",
		Settings:     map[string]string{"theme": "dark", "language": "en"},
		PasswordIDs:  []string{"test-entry-1", "test-entry-2", "test-entry-3"},
		UpdatedAt:    time.Now(),
		TotalEntries: 3,
	}
}

// GetTestBlockchainMetadata создает тестовые метаданные блокчейна
func GetTestBlockchainMetadata() *vault.BlockchainMetadata {
	return &vault.BlockchainMetadata{
		Version:      "1.0",
		UpdatedAt:    time.Now(),
		Settings:     map[string]string{"theme": "dark", "language": "en"},
		TotalEntries: 3,
	}
}

// GetTestLocalVault создает тестовое локальное хранилище
func GetTestLocalVault() *vault.LocalVault {
	vault := vault.NewLocalVault()

	// Добавляем тестовые записи
	entries := GetTestPasswordEntries()
	for _, entry := range entries {
		vault.Entries[entry.ID] = entry
	}

	// Устанавливаем метаданные
	vault.Metadata = GetTestUserMetadata()
	vault.LastSyncTime = time.Now()
	vault.IsDirty = false

	// Добавляем маппинг блокчейна
	vault.BlockchainEntries["test-entry-1"] = big.NewInt(1)
	vault.BlockchainEntries["test-entry-2"] = big.NewInt(2)
	vault.BlockchainEntries["test-entry-3"] = big.NewInt(3)

	return vault
}

// GetTestVaultConfig создает тестовую конфигурацию хранилища
func GetTestVaultConfig() *vault.VaultConfig {
	return &vault.VaultConfig{
		Argon2Time:      3,
		Argon2Memory:    64 * 1024,
		Argon2Threads:   4,
		Argon2KeyLength: 32,
	}
}

// GetTestBlockchainConfig создает тестовую конфигурацию блокчейна
func GetTestBlockchainConfig() *blockchain.BlockchainConfig {
	return &blockchain.BlockchainConfig{
		RPCEndpoint:     TestRPCEndpoint,
		ContractAddress: TestContractAddress,
		ChainID:         TestChainID,
		GasLimit:        100000,
		GasPrice:        nil,
	}
}

// GetTestEncryptedEntryBlob создает тестовый зашифрованный blob записи
func GetTestEncryptedEntryBlob() *vault.EncryptedEntryBlob {
	return &vault.EncryptedEntryBlob{
		EncryptedData: []byte("encrypted-entry-data"),
		Salt:          []byte("test-salt-16-by"),
		Nonce:         []byte("test-nonce-12"),
	}
}

// GetTestEncryptedMetadataBlob создает тестовый зашифрованный blob метаданных
func GetTestEncryptedMetadataBlob() *vault.EncryptedMetadataBlob {
	return &vault.EncryptedMetadataBlob{
		EncryptedData: []byte("encrypted-metadata-data"),
		Salt:          []byte("test-salt-16-by"),
		Nonce:         []byte("test-nonce-12"),
	}
}

// GetTestSyncStatus создает тестовый статус синхронизации
func GetTestSyncStatus() *vault.SyncStatus {
	return &vault.SyncStatus{
		LastSyncTime:   time.Now(),
		PendingChanges: map[string]string{"test-entry-1": "update"},
		FailedSyncs:    0,
		IsOnline:       true,
	}
}

// GetTestDataIDs создает тестовые ID данных
func GetTestDataIDs() []*big.Int {
	return []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
	}
}

// GetTestMasterKey создает тестовый мастер-ключ
func GetTestMasterKey() []byte {
	return []byte("test-master-key-32-bytes-long-123")
}

// GetTestSalt создает тестовую соль
func GetTestSalt() []byte {
	return []byte("test-salt-16-by")
}

// GetTestNonce создает тестовый nonce
func GetTestNonce() []byte {
	return []byte("test-nonce-12")
}

// GetTestCiphertext создает тестовый зашифрованный текст
func GetTestCiphertext() []byte {
	return []byte("test-ciphertext-data")
}

// GetTestPlaintext создает тестовый открытый текст
func GetTestPlaintext() []byte {
	return []byte("test-plaintext-data")
}

// GetTestArgon2Config создает тестовую конфигурацию Argon2
func GetTestArgon2Config() map[string]interface{} {
	return map[string]interface{}{
		"time":       uint32(3),
		"memory":     uint32(64 * 1024),
		"threads":    uint8(4),
		"key_length": uint32(32),
	}
}

// GetTestTransactionHash создает тестовый хеш транзакции
func GetTestTransactionHash() string {
	return "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
}

// GetTestBlockNumber создает тестовый номер блока
func GetTestBlockNumber() *big.Int {
	return big.NewInt(12345)
}

// GetTestGasUsed создает тестовое количество использованного газа
func GetTestGasUsed() uint64 {
	return 50000
}

// GetTestGasPrice создает тестовую цену газа
func GetTestGasPrice() *big.Int {
	return big.NewInt(20000000000) // 20 gwei
}

// GetTestNetworkID создает тестовый ID сети
func GetTestNetworkID() *big.Int {
	return big.NewInt(TestChainID)
}

// GetTestError создает тестовую ошибку
func GetTestError() error {
	return &blockchain.BlockchainError{
		Type:    "TEST_ERROR",
		Message: "Test error message",
		Code:    1001,
	}
}
