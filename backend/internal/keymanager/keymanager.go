package keymanager

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

	return &KeyManager{
		ConfigDir:     config.ConfigDir,
		SessionActive: false,
		LastActivity:  time.Now(),
	}
}

func (km *KeyManager) HasStoredKeys() bool {
	keyFilePath := km.getKeyFilePath()
	_, err := os.Stat(keyFilePath)
	return err == nil
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
