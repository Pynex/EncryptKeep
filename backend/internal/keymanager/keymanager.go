package keymanager

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"

	"encryptkeep-backend/internal/codec"
	localcrypto "encryptkeep-backend/internal/crypto"
	"encryptkeep-backend/internal/vault"

	"github.com/ethereum/go-ethereum/crypto"
)

type KeyManager struct {
	ConfigDir     string
	PrivateKey    *ecdsa.PrivateKey
	MasterKey     []byte
	SessionActive bool
	LastActivity  time.Time
	config        KeyManagerConfig
}

type KeyManagerConfig struct {
	ConfigDir      string
	SessionTimeout time.Duration
}

type StoredKeyData struct {
	EncryptedPrivateKey []byte `json:"encrypted_private_key"`
	Salt                []byte `json:"salt"`
	Nonce               []byte `json:"nonce"`
	CreatedAt           string `json:"created_at"`
	Address             string `json:"address"`
}

func NewKeyManager(config KeyManagerConfig) *KeyManager {
	if config.SessionTimeout == 0 {
		config.SessionTimeout = 60 * 4 * time.Minute
	}

	if config.ConfigDir == "" {
		config.ConfigDir = defaultConfigDir()
	}

	return &KeyManager{
		ConfigDir:     config.ConfigDir,
		SessionActive: false,
		LastActivity:  time.Now(),
		config:        config,
	}
}

// - Windows: %AppData%\encryptkeep\keys
// - Linux:   $XDG_CONFIG_HOME/encryptkeep/keys or ~/.config/encryptkeep/keys
// - macOS:   ~/Library/Application Support/encryptkeep/keys
func defaultConfigDir() string {
	if v := os.Getenv("ENCRYPTKEEP_CONFIG_DIR"); v != "" {
		return v
	}

	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		home, _ := os.UserHomeDir()
		if home == "" {
			return "encryptkeep/keys"
		}
		base = filepath.Join(home, ".config")
	}

	return filepath.Join(base, "encryptkeep", "keys")
}

func (km *KeyManager) HasStoredKeys() bool {
	keyFilePath := km.getKeyFilePath()
	_, err := os.Stat(keyFilePath)
	return err == nil
}

func (km *KeyManager) GetAddress() (string, error) {
	if !km.IsSessionActive() {
		return "", fmt.Errorf("session is not active")
	}

	publicKey := km.PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("invalid public key")
	}
	userAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	km.UpdateActivity()

	return userAddress.Hex(), nil
}

func (km *KeyManager) getKeyFilePath() string {
	return filepath.Join(km.ConfigDir, "keys.json")
}

func (km *KeyManager) saveKeyData(data StoredKeyData) error {
	if err := os.MkdirAll(km.ConfigDir, 0700); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	keyFilePath := km.getKeyFilePath()
	return os.WriteFile(keyFilePath, jsonData, 0600)
}

func (km *KeyManager) loadKeyData() (StoredKeyData, error) {
	keyFilePath := km.getKeyFilePath()

	jsonData, err := os.ReadFile(keyFilePath)
	if err != nil {
		return StoredKeyData{}, err
	}

	var data StoredKeyData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return StoredKeyData{}, err
	}

	return data, nil
}

func (km *KeyManager) IsSessionActive() bool {
	if !km.SessionActive {
		return false
	}

	if time.Since(km.LastActivity) > km.config.SessionTimeout {
		km.ClearSession()
		return false
	}

	return true
}

func (km *KeyManager) ClearSession() {
	if km.PrivateKey != nil {
		km.PrivateKey = nil
	}

	if km.MasterKey != nil {
		for i := range km.MasterKey {
			km.MasterKey[i] = 0
		}
		km.MasterKey = nil
	}

	km.SessionActive = false
}

func (km *KeyManager) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	if !km.IsSessionActive() {
		return nil, fmt.Errorf("session not active")
	}

	km.LastActivity = time.Now()
	return km.PrivateKey, nil
}

func (km *KeyManager) GetMasterKey() ([]byte, error) {
	if !km.IsSessionActive() {
		return nil, fmt.Errorf("session not active")
	}

	km.LastActivity = time.Now()
	return km.MasterKey, nil
}

func (km *KeyManager) UpdateActivity() {
	if km.SessionActive {
		km.LastActivity = time.Now()
	}
}

func (km *KeyManager) LoadFromStorage(masterPassword string) error {
	loadedData, err := km.loadKeyData()
	if err != nil {
		return err
	}

	encrypted := localcrypto.Sealed{
		Salt:       loadedData.Salt,
		Nonce:      loadedData.Nonce,
		Ciphertext: loadedData.EncryptedPrivateKey,
	}

	privateKeyHex, err := localcrypto.Open(masterPassword, codec.FromVaultConfig(vault.DefaultVaultConfig()), encrypted)
	if err != nil {
		return err
	}

	privateKey, err := crypto.HexToECDSA(string(privateKeyHex))
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid public key")
	}
	userAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	if userAddress.Hex() != loadedData.Address {
		return fmt.Errorf("addresses mismatch")
	}

	return km.startSession(privateKey, masterPassword)
}

func (km *KeyManager) InitializeFirstTime(privateKeyHex, masterPassword string) error {
	if len(privateKeyHex) != 64 {
		return fmt.Errorf("invalid private key")
	}

	if len(masterPassword) < 8 {
		return fmt.Errorf("master password should be greater or equal 8")
	}

	if !utf8.ValidString(privateKeyHex) {
		return fmt.Errorf("invalid private key hex")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid public key")
	}
	userAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	argon := codec.FromVaultConfig(vault.DefaultVaultConfig())

	sealed, err := localcrypto.Seal(masterPassword, argon, []byte(privateKeyHex))
	if err != nil {
		return err
	}

	storedData := StoredKeyData{
		EncryptedPrivateKey: sealed.Ciphertext,
		Salt:                sealed.Salt,
		Nonce:               sealed.Nonce,
		CreatedAt:           time.Now().Format(time.RFC3339),
		Address:             userAddress.Hex(),
	}

	if err := km.saveKeyData(storedData); err != nil {
		return err
	}

	return km.startSession(privateKey, masterPassword)
}

func (km *KeyManager) startSession(privateKey *ecdsa.PrivateKey, masterPassword string) error {
	dk, err := localcrypto.DeriveKey(masterPassword, nil, codec.FromVaultConfig(vault.DefaultVaultConfig()))
	if err != nil {
		return err
	}

	km.PrivateKey = privateKey
	km.MasterKey = dk.Key
	km.SessionActive = true
	km.LastActivity = time.Now()

	return nil
}
