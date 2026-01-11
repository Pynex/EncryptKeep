package vaultmanager

import (
	"context"
	"fmt"

	"encryptkeep-backend/internal/blockchain"
	"encryptkeep-backend/internal/codec"
	"encryptkeep-backend/internal/vault"
)

type VaultManager struct {
	service        blockchain.BlockchainService
	codec          *codec.Codec
	masterPassword string
}

func NewVaultManager(service blockchain.BlockchainService, masterPassword string) *VaultManager {
	return &VaultManager{
		service:        service,
		codec:          codec.NewCodec(),
		masterPassword: masterPassword,
	}
}

func (vm *VaultManager) AddEntry(ctx context.Context, v *vault.LocalVault, entry *vault.PasswordEntry) error {
	if entry == nil {
		return fmt.Errorf("entry is nil")
	}
	data, err := vm.codec.PackEntry(entry, vm.masterPassword)
	if err != nil {
		return err
	}

	if _, err := vm.service.StoreData(ctx, data); err != nil {
		return err
	}
	return vm.service.SyncVault(v)
}

func (vm *VaultManager) UpdateEntry(ctx context.Context, v *vault.LocalVault, entry *vault.PasswordEntry) error {
	if entry == nil {
		return fmt.Errorf("entry is nil")
	}
	contractID, ok := v.BlockchainEntries[entry.ID]
	if !ok {
		return fmt.Errorf("contract id not found for entry %s", entry.ID)
	}

	data, err := vm.codec.PackEntry(entry, vm.masterPassword)
	if err != nil {
		return err
	}

	if _, err := vm.service.ChangeData(ctx, contractID, data); err != nil {
		return err
	}
	return vm.service.SyncVault(v)
}

func (vm *VaultManager) DeleteEntry(ctx context.Context, v *vault.LocalVault, entryID string) error {
	contractID, ok := v.BlockchainEntries[entryID]
	if !ok {
		return fmt.Errorf("contract id not found for entry %s", entryID)
	}

	if _, err := vm.service.RemoveData(ctx, contractID); err != nil {
		return err
	}
	return vm.service.SyncVault(v)
}

func (vm *VaultManager) GetEntryFromVault(v *vault.LocalVault, entryID string) (*vault.PasswordEntry, error) {
	if entryID == "" {
		return nil, fmt.Errorf("entryID is empty")
	}
	entry, ok := v.Entries[entryID]
	if !ok {
		return nil, fmt.Errorf("entry not found: %s", entryID)
	}
	return entry, nil
}

func (vm *VaultManager) GetAllEntries(v *vault.LocalVault) []*vault.PasswordEntry {
	result := make([]*vault.PasswordEntry, 0, len(v.Entries))
	for _, e := range v.Entries {
		result = append(result, e)
	}
	return result
}

func (vm *VaultManager) StoreMetadata(ctx context.Context, meta *vault.UserMetadata) error {
	if meta == nil {
		return fmt.Errorf("metadata is nil")
	}
	data, err := vm.codec.PackMetadata(meta, vm.masterPassword)
	if err != nil {
		return err
	}
	if _, err := vm.service.StoreMetadata(ctx, data); err != nil {
		return err
	}
	return vm.service.SyncVault(&vault.LocalVault{
		Metadata: meta,
	})
}
