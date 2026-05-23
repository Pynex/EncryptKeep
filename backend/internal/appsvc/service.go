package appsvc

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"encryptkeep-backend/internal/blockchain"
	"encryptkeep-backend/internal/keymanager"
	"encryptkeep-backend/internal/vault"
	"encryptkeep-backend/internal/vaultmanager"

	"github.com/ethereum/go-ethereum/crypto"
)

// Service holds shared application state for CLI and GUI entrypoints.
type Service struct {
	mu sync.Mutex

	km *keymanager.KeyManager

	svc          blockchain.BlockchainService
	localVault   *vault.LocalVault
	vm           *vaultmanager.VaultManager
	masterPassword string
}

// NewService builds a service with default key storage paths and timeouts.
func NewService() *Service {
	return &Service{
		km: keymanager.NewKeyManager(keymanager.KeyManagerConfig{}),
	}
}

// HasStoredKeys reports whether encrypted keys exist on disk.
func (s *Service) HasStoredKeys() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.km.HasStoredKeys()
}

// IsUnlocked reports whether blockchain and vault are ready for operations.
func (s *Service) IsUnlocked() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.localVault != nil && s.vm != nil && s.svc != nil
}

// Unlock loads keys from disk and connects to the chain.
func (s *Service) Unlock(ctx context.Context, masterPassword string) error {
	originalPassword := masterPassword
	masterPassword = normalizeMasterPassword(masterPassword)
	if len(masterPassword) < 8 {
		return fmt.Errorf("master password too short (min 8)")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.resetBlockchainState()
	if err := s.km.LoadFromStorage(masterPassword); err != nil {
		// Keep compatibility with passwords that intentionally end with CR/LF.
		if originalPassword != masterPassword {
			if retryErr := s.km.LoadFromStorage(originalPassword); retryErr == nil {
				masterPassword = originalPassword
			} else {
				return fmt.Errorf("unlock local keys: %w", err)
			}
		} else {
			return fmt.Errorf("unlock local keys: %w", err)
		}
	}
	return s.connectAndUnlockVault(ctx, masterPassword)
}

// InitializeNewKeys creates encrypted keys from a hex private key and connects.
func (s *Service) InitializeNewKeys(ctx context.Context, privateKeyHex, masterPassword string) error {
	masterPassword = normalizeMasterPassword(masterPassword)
	if len(masterPassword) < 8 {
		return fmt.Errorf("master password too short (min 8)")
	}
	privateKeyHex = strings.TrimSpace(privateKeyHex)
	if len(privateKeyHex) != 64 {
		return fmt.Errorf("invalid private key length")
	}
	if _, err := hex.DecodeString(privateKeyHex); err != nil {
		return fmt.Errorf("invalid private key hex: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.resetBlockchainState()
	if s.km.HasStoredKeys() {
		return fmt.Errorf("keys already exist; use unlock instead")
	}
	if err := s.km.InitializeFirstTime(privateKeyHex, masterPassword); err != nil {
		return err
	}
	return s.connectAndUnlockVault(ctx, masterPassword)
}

func (s *Service) resetBlockchainState() {
	if s.svc != nil {
		_ = s.svc.Disconnect()
	}
	s.svc = nil
	s.localVault = nil
	s.vm = nil
	s.masterPassword = ""
}

func (s *Service) connectAndUnlockVault(ctx context.Context, masterPassword string) error {
	privKey, err := s.km.GetPrivateKey()
	if err != nil {
		return err
	}
	privHex := hex.EncodeToString(crypto.FromECDSA(privKey))

	svc := blockchain.NewBlockchainService(blockchain.GetDefaultConfig())
	if err := svc.Connect(); err != nil {
		return fmt.Errorf("blockchain connect: %w", err)
	}
	if _, err := svc.StartSession(privHex, masterPassword); err != nil {
		_ = svc.Disconnect()
		return fmt.Errorf("start session: %w", err)
	}

	localVault := vault.NewLocalVault()
	if err := svc.SyncVault(localVault); err != nil {
		_ = svc.Disconnect()
		return fmt.Errorf("sync vault: %w", err)
	}

	s.svc = svc
	s.localVault = localVault
	s.masterPassword = masterPassword
	s.vm = vaultmanager.NewVaultManager(svc, masterPassword)
	_ = ctx // reserved for cancellable ops later
	return nil
}

func normalizeMasterPassword(masterPassword string) string {
	return strings.TrimRight(masterPassword, "\r\n")
}

// Lock disconnects from the chain and clears the key session.
func (s *Service) Lock() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetBlockchainState()
	s.km.ClearSession()
}

// ResetStoredKeys clears the active session and deletes local encrypted keys.
func (s *Service) ResetStoredKeys() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetBlockchainState()
	if err := s.km.DeleteStoredKeys(); err != nil {
		return fmt.Errorf("delete stored keys: %w", err)
	}
	return nil
}

// ListEntries returns all password entries in the local vault.
func (s *Service) ListEntries() ([]*vault.PasswordEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vm == nil || s.localVault == nil {
		return nil, fmt.Errorf("not unlocked")
	}
	return s.vm.GetAllEntries(s.localVault), nil
}

// GetEntry returns one entry by id.
func (s *Service) GetEntry(id string) (*vault.PasswordEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vm == nil || s.localVault == nil {
		return nil, fmt.Errorf("not unlocked")
	}
	return s.vm.GetEntryFromVault(s.localVault, id)
}

// AddEntry creates and syncs a new entry.
func (s *Service) AddEntry(ctx context.Context, title, username, password, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vm == nil || s.localVault == nil {
		return fmt.Errorf("not unlocked")
	}
	entry := vault.NewPasswordEntry(title, username, password)
	entry.URL = url
	return s.vm.AddEntry(ctx, s.localVault, entry)
}

// UpdateEntry updates an existing entry and syncs.
func (s *Service) UpdateEntry(ctx context.Context, entry *vault.PasswordEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vm == nil || s.localVault == nil {
		return fmt.Errorf("not unlocked")
	}
	if entry == nil {
		return fmt.Errorf("entry is nil")
	}
	entry.UpdatedAt = time.Now()
	return s.vm.UpdateEntry(ctx, s.localVault, entry)
}

// DeleteEntry removes an entry by id and syncs.
func (s *Service) DeleteEntry(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vm == nil || s.localVault == nil {
		return fmt.Errorf("not unlocked")
	}
	return s.vm.DeleteEntry(ctx, s.localVault, id)
}

// Sync pulls remote state into the local vault.
func (s *Service) Sync(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.svc == nil || s.localVault == nil {
		return fmt.Errorf("not unlocked")
	}
	if err := s.svc.SyncVault(s.localVault); err != nil {
		return err
	}
	_ = ctx
	return nil
}

// VaultStats returns entry count and last sync time after a successful unlock/sync.
func (s *Service) VaultStats() (entryCount int, lastSync time.Time, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.localVault == nil {
		return 0, time.Time{}, fmt.Errorf("not unlocked")
	}
	return len(s.localVault.Entries), s.localVault.LastSyncTime, nil
}
