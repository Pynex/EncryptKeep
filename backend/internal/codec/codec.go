package codec

import (
	"encoding/json"
	// "math/big"
	"encryptkeep-backend/internal/crypto"
	"encryptkeep-backend/internal/vault"
	"errors"
	"fmt"
)

type Codec struct {
	argon2Config crypto.Argon2Config
}

func NewCodec() *Codec {
	return &Codec{
		argon2Config: FromVaultConfig(vault.DefaultVaultConfig()),
	}
}

func NewCodecWithConfig(cfg crypto.Argon2Config) *Codec {
	return &Codec{
		argon2Config: cfg,
	}
}

func NewCodecWithVaultConfig(vaultCfg *vault.VaultConfig) *Codec {
	return &Codec{
		argon2Config: FromVaultConfig(vaultCfg),
	}
}

func (c *Codec) PackEntry(passwordEntry *vault.PasswordEntry, masterPassword string) ([]byte, error) {
	if passwordEntry == nil {
		return nil, errors.New("passwordEntry cannot be empty")
	}

	if len(masterPassword) < 8 { //min 8
		return nil, errors.New("master password cannot be empty")
	}

	plaintext, err := json.Marshal(passwordEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry: %w", err)
	}

	sealed, err := crypto.Seal(masterPassword, c.argon2Config, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to seal plain text: %w", err)
	}

	encryptedBlob := &vault.EncryptedEntryBlob{
		EncryptedData: sealed.Ciphertext,
		Salt:          sealed.Salt,
		Nonce:         sealed.Nonce,
	}

	packedData, err := json.Marshal(encryptedBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to convert marshal sealed data: %w", err)
	}

	return packedData, nil
}

func (c *Codec) PackMetadata(metadata *vault.UserMetadata, masterPassword string) ([]byte, error) {
	if metadata == nil {
		return nil, errors.New("metadata cannot be empty")
	}

	if len(masterPassword) < 8 {
		return nil, errors.New("master password cannot be empty")
	}

	blockchainMetaData := vault.BlockchainMetadata{
		Version:      metadata.Version,
		UpdatedAt:    metadata.UpdatedAt,
		Settings:     metadata.Settings,
		TotalEntries: metadata.TotalEntries,
	}

	plaintext, err := json.Marshal(blockchainMetaData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert marshal metadata: %w", err)
	}

	sealed, err := crypto.Seal(masterPassword, c.argon2Config, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to seal data: %w", err)
	}

	encryptedBlob := &vault.EncryptedMetadataBlob{
		EncryptedData: sealed.Ciphertext,
		Salt:          sealed.Salt,
		Nonce:         sealed.Nonce,
	}

	packedData, err := json.Marshal(encryptedBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to convert marshal sealed data: %w", err)
	}

	return packedData, nil
}

func (c *Codec) UnpackEntry(encryptedData []byte, masterPassword string) (*vault.PasswordEntry, error) {
	if len(masterPassword) < 8 {
		return nil, errors.New("master password cannot be empty")
	}

	length := len(encryptedData)
	if length == 0 {
		return nil, errors.New("encrypted data cannot be empty")
	}

	var encryptedBlob vault.EncryptedEntryBlob
	if err := json.Unmarshal(encryptedData, &encryptedBlob); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata^ %w", err)
	}
	sealed := crypto.Sealed{
		Salt:       encryptedBlob.Salt,
		Nonce:      encryptedBlob.Nonce,
		Ciphertext: encryptedBlob.EncryptedData,
	}

	plaintext, err := crypto.Open(masterPassword, c.argon2Config, sealed)
	if err != nil {
		return nil, fmt.Errorf("failed to open plaintext: %w", err)
	}

	var entry vault.PasswordEntry
	if err := json.Unmarshal(plaintext, &entry); err != nil {
		return nil, fmt.Errorf("failed unmarshal plaintext: %w", err)
	}

	return &entry, nil
}

func FromVaultConfig(vaultCfg *vault.VaultConfig) crypto.Argon2Config {
	return crypto.Argon2Config{
		Time:      vaultCfg.Argon2Time,
		Memory:    vaultCfg.Argon2Memory,
		Threads:   vaultCfg.Argon2Threads,
		KeyLength: vaultCfg.Argon2KeyLength,
	}
}

func (c *Codec) UnpackMetadata(encryptedData []byte, masterPassword string) (*vault.UserMetadata, error) {
	if len(masterPassword) < 8 {
		return nil, errors.New("master password cannot be empty")
	}

	length := len(encryptedData)
	if length == 0 {
		return nil, errors.New("encrypted data cannot be empty")
	}

	var encryptedBlob vault.EncryptedMetadataBlob
	if err := json.Unmarshal(encryptedData, &encryptedBlob); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata^ %w", err)
	}
	sealed := crypto.Sealed{
		Salt:       encryptedBlob.Salt,
		Nonce:      encryptedBlob.Nonce,
		Ciphertext: encryptedBlob.EncryptedData,
	}

	plaintext, err := crypto.Open(masterPassword, c.argon2Config, sealed)
	if err != nil {
		return nil, fmt.Errorf("failed to open plaintext: %w", err)
	}

	var blockchainMeta vault.BlockchainMetadata
	if err := json.Unmarshal(plaintext, &blockchainMeta); err != nil {
		return nil, fmt.Errorf("failed unmarshal plaintext: %w", err)
	}

	metadata := &vault.UserMetadata{
		Version:      blockchainMeta.Version,
		Settings:     blockchainMeta.Settings,
		PasswordIDs:  []string{},
		UpdatedAt:    blockchainMeta.UpdatedAt,
		TotalEntries: blockchainMeta.TotalEntries,
	}

	return metadata, nil
}
