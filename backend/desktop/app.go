package main

import (
	"context"
	"time"

	"encryptkeep-backend/internal/appsvc"
	"encryptkeep-backend/internal/vault"
)

// App exposes EncryptKeep operations to the Wails frontend.
type App struct {
	ctx context.Context
	svc *appsvc.Service
}

// VaultStats is returned to the UI after unlock or sync.
type VaultStats struct {
	EntryCount int    `json:"entryCount"`
	LastSync   string `json:"lastSync"`
}

// NewApp creates the desktop application API.
func NewApp() *App {
	return &App{svc: appsvc.NewService()}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	a.svc.Lock()
}

// HasStoredKeys reports whether encrypted keys exist on disk.
func (a *App) HasStoredKeys() bool {
	return a.svc.HasStoredKeys()
}

// IsUnlocked reports whether the vault session is active.
func (a *App) IsUnlocked() bool {
	return a.svc.IsUnlocked()
}

// Unlock loads keys and syncs the vault with the master password.
func (a *App) Unlock(masterPassword string) error {
	return a.svc.Unlock(a.ctx, masterPassword)
}

// InitializeNewKeys stores a new private key and unlocks the vault.
func (a *App) InitializeNewKeys(privateKeyHex, masterPassword string) error {
	return a.svc.InitializeNewKeys(a.ctx, privateKeyHex, masterPassword)
}

// Lock ends the session and disconnects from the chain.
func (a *App) Lock() {
	a.svc.Lock()
}

// ListEntries returns all password entries.
func (a *App) ListEntries() ([]*vault.PasswordEntry, error) {
	return a.svc.ListEntries()
}

// GetEntry returns one entry by id.
func (a *App) GetEntry(id string) (*vault.PasswordEntry, error) {
	return a.svc.GetEntry(id)
}

// AddEntry creates and syncs a new entry.
func (a *App) AddEntry(title, username, password, url string) error {
	return a.svc.AddEntry(a.ctx, title, username, password, url)
}

// UpdateEntry updates an existing entry and syncs.
func (a *App) UpdateEntry(entry vault.PasswordEntry) error {
	e := entry
	return a.svc.UpdateEntry(a.ctx, &e)
}

// DeleteEntry removes an entry by id.
func (a *App) DeleteEntry(id string) error {
	return a.svc.DeleteEntry(a.ctx, id)
}

// Sync pulls the latest vault state from the chain.
func (a *App) Sync() error {
	return a.svc.Sync(a.ctx)
}

// GetVaultStats returns entry count and last sync time.
func (a *App) GetVaultStats() (VaultStats, error) {
	n, last, err := a.svc.VaultStats()
	if err != nil {
		return VaultStats{}, err
	}
	return VaultStats{
		EntryCount: n,
		LastSync:   last.Format(time.RFC3339),
	}, nil
}
