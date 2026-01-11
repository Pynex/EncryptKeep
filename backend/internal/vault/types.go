package vault

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"
)

type uint256 = *big.Int

type PasswordEntry struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsFavorite bool      `json:"is_favorite"`
}

type BlockchainEntry struct {
	ContractID uint256        `json:"contract_id"`
	Entry      *PasswordEntry `json:"entry"`
}

type BlockchainMetadata struct {
	Version      string            `json:"version"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Settings     map[string]string `json:"settings"`
	TotalEntries int               `json:"total_entries"`
}

type EncryptedEntryBlob struct {
	EncryptedData []byte `json:"encrypted_data"`
	Salt          []byte `json:"salt"`
	Nonce         []byte `json:"nonce"`
}

type EncryptedMetadataBlob struct {
	EncryptedData []byte `json:"encrypted_data"`
	Salt          []byte `json:"salt"`
	Nonce         []byte `json:"nonce"`
}

type UserMetadata struct {
	Version      string            `json:"version"`
	Settings     map[string]string `json:"settings"`
	PasswordIDs  []string          `json:"password_ids"`
	UpdatedAt    time.Time         `json:"updated_at"`
	TotalEntries int               `json:"total_entries"`
}

type LocalVault struct {
	Entries      map[string]*PasswordEntry `json:"entries"`
	Metadata     *UserMetadata             `json:"metadata"`
	LastSyncTime time.Time                 `json:"last_sync_time"`
	IsDirty      bool                      `json:"is_dirty"` // unsaved changes

	BlockchainEntries map[string]uint256 `json:"blockchain_entries"`
}

type MasterKey struct {
	Key  []byte
	Salt []byte
}

type VaultConfig struct {
	Argon2Time      uint32
	Argon2Memory    uint32
	Argon2Threads   uint8
	Argon2KeyLength uint32
}

type SyncStatus struct {
	LastSyncTime   time.Time         `json:"last_sync_time"`
	PendingChanges map[string]string `json:"pending_changes"` // ID -> "add"/"update"/"delete"
	FailedSyncs    int               `json:"failed_syncs"`
	IsOnline       bool              `json:"is_online"`
}

func DefaultVaultConfig() *VaultConfig {
	return &VaultConfig{
		Argon2Time:      3, 
		Argon2Memory:    64 * 1024,
		Argon2Threads:   4,
		Argon2KeyLength: 32,
	}
}

func NewLocalVault() *LocalVault {
	return &LocalVault{
		Entries: make(map[string]*PasswordEntry),
		Metadata: &UserMetadata{
			Version:      "1.0",
			Settings:     make(map[string]string),
			PasswordIDs:  []string{},
			UpdatedAt:    time.Now(),
			TotalEntries: 0,
		},
		LastSyncTime:      time.Now(),
		IsDirty:           false,
		BlockchainEntries: make(map[string]uint256),
	}
}

func NewPasswordEntry(title, username, password string) *PasswordEntry {
	now := time.Now()
	return &PasswordEntry{
		ID:         generateID(),
		Title:      title,
		Username:   username,
		Password:   password,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsFavorite: false,
	}
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format("20060102150405999")
	}

	return hex.EncodeToString(b)
}
