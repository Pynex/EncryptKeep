package vault_test

import (
	"math/big"
	"testing"
	"time"

	"encryptkeep-backend/internal/vault"
)

// TestDefaultVaultConfig тестирует создание конфигурации по умолчанию
func TestDefaultVaultConfig(t *testing.T) {
	config := vault.DefaultVaultConfig()

	if config == nil {
		t.Fatal("DefaultVaultConfig should not return nil")
	}

	if config.Argon2Time == 0 {
		t.Error("Argon2Time should not be zero")
	}

	if config.Argon2Memory == 0 {
		t.Error("Argon2Memory should not be zero")
	}

	if config.Argon2Threads == 0 {
		t.Error("Argon2Threads should not be zero")
	}

	if config.Argon2KeyLength == 0 {
		t.Error("Argon2KeyLength should not be zero")
	}

	// Проверяем разумные значения по умолчанию
	if config.Argon2Time < 1 {
		t.Error("Argon2Time should be at least 1")
	}

	if config.Argon2Memory < 1024 {
		t.Error("Argon2Memory should be at least 1KB")
	}

	if config.Argon2Threads < 1 {
		t.Error("Argon2Threads should be at least 1")
	}

	if config.Argon2KeyLength < 16 {
		t.Error("Argon2KeyLength should be at least 16 bytes")
	}
}

// TestNewLocalVault тестирует создание нового локального хранилища
func TestNewLocalVault(t *testing.T) {
	vault := vault.NewLocalVault()

	if vault == nil {
		t.Fatal("NewLocalVault should not return nil")
	}

	if vault.Entries == nil {
		t.Error("Entries should not be nil")
	}

	if vault.Metadata == nil {
		t.Error("Metadata should not be nil")
	}

	if vault.BlockchainEntries == nil {
		t.Error("BlockchainEntries should not be nil")
	}

	// Проверяем метаданные по умолчанию
	if vault.Metadata.Version == "" {
		t.Error("Default version should not be empty")
	}

	if vault.Metadata.Settings == nil {
		t.Error("Default settings should not be nil")
	}

	if vault.Metadata.PasswordIDs == nil {
		t.Error("Default PasswordIDs should not be nil")
	}

	if vault.Metadata.TotalEntries != 0 {
		t.Error("Default TotalEntries should be 0")
	}
}

// TestNewPasswordEntry тестирует создание новой записи пароля
func TestNewPasswordEntry(t *testing.T) {
	title := "Test Site"
	username := "testuser"
	password := "testpassword"

	entry := vault.NewPasswordEntry(title, username, password)

	if entry == nil {
		t.Fatal("NewPasswordEntry should not return nil")
	}

	if entry.ID == "" {
		t.Error("ID should not be empty")
	}

	if entry.Title != title {
		t.Errorf("Title mismatch: got %s, want %s", entry.Title, title)
	}

	if entry.Username != username {
		t.Errorf("Username mismatch: got %s, want %s", entry.Username, username)
	}

	if entry.Password != password {
		t.Errorf("Password mismatch: got %s, want %s", entry.Password, password)
	}

	if entry.URL != "" {
		t.Error("URL should be empty by default")
	}

	if entry.IsFavorite {
		t.Error("IsFavorite should be false by default")
	}

	// Проверяем, что временные метки установлены
	if entry.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if entry.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	// Проверяем, что CreatedAt и UpdatedAt одинаковые при создании
	if !entry.CreatedAt.Equal(entry.UpdatedAt) {
		t.Error("CreatedAt and UpdatedAt should be equal at creation time")
	}
}

// TestPasswordEntryIDGeneration тестирует генерацию уникальных ID
func TestPasswordEntryIDGeneration(t *testing.T) {
	// Создаем несколько записей
	entry1 := vault.NewPasswordEntry("Site 1", "user1", "pass1")
	entry2 := vault.NewPasswordEntry("Site 2", "user2", "pass2")
	entry3 := vault.NewPasswordEntry("Site 3", "user3", "pass3")

	// Проверяем, что ID уникальны
	if entry1.ID == entry2.ID {
		t.Error("Entry IDs should be unique")
	}

	if entry1.ID == entry3.ID {
		t.Error("Entry IDs should be unique")
	}

	if entry2.ID == entry3.ID {
		t.Error("Entry IDs should be unique")
	}

	// Проверяем, что ID не пустые
	if entry1.ID == "" {
		t.Error("Entry ID should not be empty")
	}

	if entry2.ID == "" {
		t.Error("Entry ID should not be empty")
	}

	if entry3.ID == "" {
		t.Error("Entry ID should not be empty")
	}
}

// TestPasswordEntryTimestamps тестирует временные метки записей
func TestPasswordEntryTimestamps(t *testing.T) {
	beforeCreation := time.Now()

	entry := vault.NewPasswordEntry("Test Site", "testuser", "testpassword")

	afterCreation := time.Now()

	// Проверяем, что временные метки находятся в разумном диапазоне
	if entry.CreatedAt.Before(beforeCreation) {
		t.Error("CreatedAt should not be before creation time")
	}

	if entry.CreatedAt.After(afterCreation) {
		t.Error("CreatedAt should not be after creation time")
	}

	if entry.UpdatedAt.Before(beforeCreation) {
		t.Error("UpdatedAt should not be before creation time")
	}

	if entry.UpdatedAt.After(afterCreation) {
		t.Error("UpdatedAt should not be after creation time")
	}
}

// TestPasswordEntryStruct тестирует структуру PasswordEntry
func TestPasswordEntryStruct(t *testing.T) {
	entry := &vault.PasswordEntry{
		ID:         "test-id",
		Title:      "Test Title",
		Username:   "testuser",
		Password:   "testpassword",
		URL:        "https://example.com",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsFavorite: true,
	}

	// Проверяем, что все поля установлены
	if entry.ID != "test-id" {
		t.Errorf("ID mismatch: got %s, want test-id", entry.ID)
	}

	if entry.Title != "Test Title" {
		t.Errorf("Title mismatch: got %s, want Test Title", entry.Title)
	}

	if entry.Username != "testuser" {
		t.Errorf("Username mismatch: got %s, want testuser", entry.Username)
	}

	if entry.Password != "testpassword" {
		t.Errorf("Password mismatch: got %s, want testpassword", entry.Password)
	}

	if entry.URL != "https://example.com" {
		t.Errorf("URL mismatch: got %s, want https://example.com", entry.URL)
	}

	if !entry.IsFavorite {
		t.Error("IsFavorite should be true")
	}
}

// TestUserMetadataStruct тестирует структуру UserMetadata
func TestUserMetadataStruct(t *testing.T) {
	metadata := &vault.UserMetadata{
		Version:      "1.0",
		Settings:     map[string]string{"theme": "dark"},
		PasswordIDs:  []string{"id1", "id2"},
		UpdatedAt:    time.Now(),
		TotalEntries: 2,
	}

	// Проверяем, что все поля установлены
	if metadata.Version != "1.0" {
		t.Errorf("Version mismatch: got %s, want 1.0", metadata.Version)
	}

	if metadata.Settings["theme"] != "dark" {
		t.Errorf("Settings theme mismatch: got %s, want dark", metadata.Settings["theme"])
	}

	if len(metadata.PasswordIDs) != 2 {
		t.Errorf("PasswordIDs length mismatch: got %d, want 2", len(metadata.PasswordIDs))
	}

	if metadata.PasswordIDs[0] != "id1" {
		t.Errorf("PasswordIDs[0] mismatch: got %s, want id1", metadata.PasswordIDs[0])
	}

	if metadata.PasswordIDs[1] != "id2" {
		t.Errorf("PasswordIDs[1] mismatch: got %s, want id2", metadata.PasswordIDs[1])
	}

	if metadata.TotalEntries != 2 {
		t.Errorf("TotalEntries mismatch: got %d, want 2", metadata.TotalEntries)
	}
}

// TestBlockchainMetadataStruct тестирует структуру BlockchainMetadata
func TestBlockchainMetadataStruct(t *testing.T) {
	metadata := &vault.BlockchainMetadata{
		Version:      "1.0",
		UpdatedAt:    time.Now(),
		Settings:     map[string]string{"theme": "dark"},
		TotalEntries: 5,
	}

	// Проверяем, что все поля установлены
	if metadata.Version != "1.0" {
		t.Errorf("Version mismatch: got %s, want 1.0", metadata.Version)
	}

	if metadata.Settings["theme"] != "dark" {
		t.Errorf("Settings theme mismatch: got %s, want dark", metadata.Settings["theme"])
	}

	if metadata.TotalEntries != 5 {
		t.Errorf("TotalEntries mismatch: got %d, want 5", metadata.TotalEntries)
	}
}

// TestLocalVaultStruct тестирует структуру LocalVault
func TestLocalVaultStruct(t *testing.T) {
	vault := &vault.LocalVault{
		Entries:           make(map[string]*vault.PasswordEntry),
		Metadata:          &vault.UserMetadata{Version: "1.0"},
		LastSyncTime:      time.Now(),
		IsDirty:           true,
		BlockchainEntries: make(map[string]*big.Int),
	}

	// Проверяем, что все поля установлены
	if vault.Entries == nil {
		t.Error("Entries should not be nil")
	}

	if vault.Metadata == nil {
		t.Error("Metadata should not be nil")
	}

	if vault.BlockchainEntries == nil {
		t.Error("BlockchainEntries should not be nil")
	}

	if vault.Metadata.Version != "1.0" {
		t.Errorf("Metadata Version mismatch: got %s, want 1.0", vault.Metadata.Version)
	}

	if !vault.IsDirty {
		t.Error("IsDirty should be true")
	}
}

// TestVaultConfigStruct тестирует структуру VaultConfig
func TestVaultConfigStruct(t *testing.T) {
	config := &vault.VaultConfig{
		Argon2Time:      3,
		Argon2Memory:    64 * 1024,
		Argon2Threads:   4,
		Argon2KeyLength: 32,
	}

	// Проверяем, что все поля установлены
	if config.Argon2Time != 3 {
		t.Errorf("Argon2Time mismatch: got %d, want 3", config.Argon2Time)
	}

	if config.Argon2Memory != 64*1024 {
		t.Errorf("Argon2Memory mismatch: got %d, want %d", config.Argon2Memory, 64*1024)
	}

	if config.Argon2Threads != 4 {
		t.Errorf("Argon2Threads mismatch: got %d, want 4", config.Argon2Threads)
	}

	if config.Argon2KeyLength != 32 {
		t.Errorf("Argon2KeyLength mismatch: got %d, want 32", config.Argon2KeyLength)
	}
}
