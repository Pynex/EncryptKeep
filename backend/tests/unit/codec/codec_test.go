package codec_test

import (
	"testing"
	"time"

	"encryptkeep-backend/internal/codec"
	"encryptkeep-backend/internal/crypto"
	"encryptkeep-backend/internal/vault"
)

// TestNewCodec тестирует создание нового Codec
func TestNewCodec(t *testing.T) {
	codec := codec.NewCodec()

	if codec == nil {
		t.Fatal("Codec should not be nil")
	}

	// Проверяем, что Codec создан
	// (не можем проверить приватные поля напрямую)
	if codec == nil {
		t.Error("Codec should not be nil")
	}
}

// TestNewCodecWithConfig тестирует создание Codec с кастомным конфигом
func TestNewCodecWithConfig(t *testing.T) {
	customConfig := crypto.Argon2Config{
		Time:      5,
		Memory:    128 * 1024,
		Threads:   8,
		KeyLength: 32,
	}

	codec := codec.NewCodecWithConfig(customConfig)

	if codec == nil {
		t.Fatal("Codec should not be nil")
	}

	// Проверяем, что Codec создан с кастомным конфигом
	// (не можем проверить приватные поля напрямую)
	if codec == nil {
		t.Error("Codec should not be nil")
	}
}

// TestNewCodecWithVaultConfig тестирует создание Codec с VaultConfig
func TestNewCodecWithVaultConfig(t *testing.T) {
	vaultConfig := &vault.VaultConfig{
		Argon2Time:      4,
		Argon2Memory:    96 * 1024,
		Argon2Threads:   6,
		Argon2KeyLength: 32,
	}

	codec := codec.NewCodecWithVaultConfig(vaultConfig)

	if codec == nil {
		t.Fatal("Codec should not be nil")
	}

	// Проверяем, что Codec создан с VaultConfig
	// (не можем проверить приватные поля напрямую)
	if codec == nil {
		t.Error("Codec should not be nil")
	}
}

// TestPackEntry тестирует упаковку PasswordEntry
func TestPackEntry(t *testing.T) {
	codec := codec.NewCodec()

	// Создаем тестовую запись
	entry := &vault.PasswordEntry{
		ID:         "test-id",
		Title:      "Test Title",
		Username:   "testuser",
		Password:   "testpassword",
		URL:        "https://example.com",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsFavorite: false,
	}

	masterPassword := "test-master-password-123"

	// Тест с nil entry
	_, err := codec.PackEntry(nil, masterPassword)
	if err == nil {
		t.Error("Should return error for nil entry")
	}

	// Тест с коротким паролем
	_, err = codec.PackEntry(entry, "short")
	if err == nil {
		t.Error("Should return error for short password")
	}

	// Тест с валидными данными
	packedData, err := codec.PackEntry(entry, masterPassword)
	if err != nil {
		t.Errorf("Should not return error for valid data: %v", err)
	}

	if packedData == nil {
		t.Error("Packed data should not be nil")
	}

	if len(packedData) == 0 {
		t.Error("Packed data should not be empty")
	}
}

// TestUnpackEntry тестирует распаковку PasswordEntry
func TestUnpackEntry(t *testing.T) {
	codec := codec.NewCodec()

	// Создаем тестовую запись
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

	masterPassword := "test-master-password-123"

	// Сначала упаковываем
	packedData, err := codec.PackEntry(entry, masterPassword)
	if err != nil {
		t.Fatalf("PackEntry failed: %v", err)
	}

	// Тест с пустыми данными
	_, err = codec.UnpackEntry([]byte{}, masterPassword)
	if err == nil {
		t.Error("Should return error for empty data")
	}

	// Тест с коротким паролем
	_, err = codec.UnpackEntry(packedData, "short")
	if err == nil {
		t.Error("Should return error for short password")
	}

	// Тест с неправильным паролем
	_, err = codec.UnpackEntry(packedData, "wrong-password")
	if err == nil {
		t.Error("Should return error for wrong password")
	}

	// Тест с правильным паролем
	unpackedEntry, err := codec.UnpackEntry(packedData, masterPassword)
	if err != nil {
		t.Errorf("Should not return error for correct password: %v", err)
	}

	if unpackedEntry == nil {
		t.Error("Unpacked entry should not be nil")
	}

	// Проверяем, что данные восстановлены правильно
	if unpackedEntry.ID != entry.ID {
		t.Errorf("ID mismatch: got %s, want %s", unpackedEntry.ID, entry.ID)
	}
	if unpackedEntry.Title != entry.Title {
		t.Errorf("Title mismatch: got %s, want %s", unpackedEntry.Title, entry.Title)
	}
	if unpackedEntry.Username != entry.Username {
		t.Errorf("Username mismatch: got %s, want %s", unpackedEntry.Username, entry.Username)
	}
	if unpackedEntry.Password != entry.Password {
		t.Errorf("Password mismatch: got %s, want %s", unpackedEntry.Password, entry.Password)
	}
	if unpackedEntry.URL != entry.URL {
		t.Errorf("URL mismatch: got %s, want %s", unpackedEntry.URL, entry.URL)
	}
	if unpackedEntry.IsFavorite != entry.IsFavorite {
		t.Errorf("IsFavorite mismatch: got %t, want %t", unpackedEntry.IsFavorite, entry.IsFavorite)
	}
}

// TestPackMetadata тестирует упаковку UserMetadata
func TestPackMetadata(t *testing.T) {
	codec := codec.NewCodec()

	// Создаем тестовые метаданные
	metadata := &vault.UserMetadata{
		Version:      "1.0",
		Settings:     map[string]string{"theme": "dark", "language": "en"},
		PasswordIDs:  []string{"id1", "id2", "id3"},
		UpdatedAt:    time.Now(),
		TotalEntries: 3,
	}

	masterPassword := "test-master-password-123"

	// Тест с nil metadata
	_, err := codec.PackMetadata(nil, masterPassword)
	if err == nil {
		t.Error("Should return error for nil metadata")
	}

	// Тест с коротким паролем
	_, err = codec.PackMetadata(metadata, "short")
	if err == nil {
		t.Error("Should return error for short password")
	}

	// Тест с валидными данными
	packedData, err := codec.PackMetadata(metadata, masterPassword)
	if err != nil {
		t.Errorf("Should not return error for valid data: %v", err)
	}

	if packedData == nil {
		t.Error("Packed data should not be nil")
	}

	if len(packedData) == 0 {
		t.Error("Packed data should not be empty")
	}
}

// TestUnpackMetadata тестирует распаковку UserMetadata
func TestUnpackMetadata(t *testing.T) {
	codec := codec.NewCodec()

	// Создаем тестовые метаданные
	metadata := &vault.UserMetadata{
		Version:      "1.0",
		Settings:     map[string]string{"theme": "dark", "language": "en"},
		PasswordIDs:  []string{"id1", "id2", "id3"},
		UpdatedAt:    time.Now(),
		TotalEntries: 3,
	}

	masterPassword := "test-master-password-123"

	// Сначала упаковываем
	packedData, err := codec.PackMetadata(metadata, masterPassword)
	if err != nil {
		t.Fatalf("PackMetadata failed: %v", err)
	}

	// Тест с пустыми данными
	_, err = codec.UnpackMetadata([]byte{}, masterPassword)
	if err == nil {
		t.Error("Should return error for empty data")
	}

	// Тест с коротким паролем
	_, err = codec.UnpackMetadata(packedData, "short")
	if err == nil {
		t.Error("Should return error for short password")
	}

	// Тест с неправильным паролем
	_, err = codec.UnpackMetadata(packedData, "wrong-password")
	if err == nil {
		t.Error("Should return error for wrong password")
	}

	// Тест с правильным паролем
	unpackedMetadata, err := codec.UnpackMetadata(packedData, masterPassword)
	if err != nil {
		t.Errorf("Should not return error for correct password: %v", err)
	}

	if unpackedMetadata == nil {
		t.Error("Unpacked metadata should not be nil")
	}

	// Проверяем, что данные восстановлены правильно
	if unpackedMetadata.Version != metadata.Version {
		t.Errorf("Version mismatch: got %s, want %s", unpackedMetadata.Version, metadata.Version)
	}
	if unpackedMetadata.TotalEntries != metadata.TotalEntries {
		t.Errorf("TotalEntries mismatch: got %d, want %d", unpackedMetadata.TotalEntries, metadata.TotalEntries)
	}

	// Проверяем настройки
	if len(unpackedMetadata.Settings) != len(metadata.Settings) {
		t.Errorf("Settings length mismatch: got %d, want %d", len(unpackedMetadata.Settings), len(metadata.Settings))
	}

	for key, value := range metadata.Settings {
		if unpackedMetadata.Settings[key] != value {
			t.Errorf("Setting %s mismatch: got %s, want %s", key, unpackedMetadata.Settings[key], value)
		}
	}

	// PasswordIDs должны быть пустыми (не сохраняются в BlockchainMetadata)
	if len(unpackedMetadata.PasswordIDs) != 0 {
		t.Errorf("PasswordIDs should be empty, got %d items", len(unpackedMetadata.PasswordIDs))
	}
}

// TestFromVaultConfig тестирует преобразование VaultConfig в Argon2Config
func TestFromVaultConfig(t *testing.T) {
	vaultConfig := &vault.VaultConfig{
		Argon2Time:      5,
		Argon2Memory:    128 * 1024,
		Argon2Threads:   8,
		Argon2KeyLength: 32,
	}

	argon2Config := codec.FromVaultConfig(vaultConfig)

	if argon2Config.Time != vaultConfig.Argon2Time {
		t.Errorf("Time mismatch: got %d, want %d", argon2Config.Time, vaultConfig.Argon2Time)
	}
	if argon2Config.Memory != vaultConfig.Argon2Memory {
		t.Errorf("Memory mismatch: got %d, want %d", argon2Config.Memory, vaultConfig.Argon2Memory)
	}
	if argon2Config.Threads != vaultConfig.Argon2Threads {
		t.Errorf("Threads mismatch: got %d, want %d", argon2Config.Threads, vaultConfig.Argon2Threads)
	}
	if argon2Config.KeyLength != vaultConfig.Argon2KeyLength {
		t.Errorf("KeyLength mismatch: got %d, want %d", argon2Config.KeyLength, vaultConfig.Argon2KeyLength)
	}
}

// TestPackUnpackRoundTrip тестирует полный цикл упаковки-распаковки
func TestPackUnpackRoundTrip(t *testing.T) {
	codec := codec.NewCodec()

	// Тестируем PasswordEntry
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

	masterPassword := "test-master-password-123"

	// Упаковываем
	packedData, err := codec.PackEntry(entry, masterPassword)
	if err != nil {
		t.Fatalf("PackEntry failed: %v", err)
	}

	// Распаковываем
	unpackedEntry, err := codec.UnpackEntry(packedData, masterPassword)
	if err != nil {
		t.Fatalf("UnpackEntry failed: %v", err)
	}

	// Проверяем, что данные идентичны
	if unpackedEntry.ID != entry.ID {
		t.Error("ID should be identical after round trip")
	}
	if unpackedEntry.Title != entry.Title {
		t.Error("Title should be identical after round trip")
	}
	if unpackedEntry.Username != entry.Username {
		t.Error("Username should be identical after round trip")
	}
	if unpackedEntry.Password != entry.Password {
		t.Error("Password should be identical after round trip")
	}
	if unpackedEntry.URL != entry.URL {
		t.Error("URL should be identical after round trip")
	}
	if unpackedEntry.IsFavorite != entry.IsFavorite {
		t.Error("IsFavorite should be identical after round trip")
	}
}

// TestPackUnpackMetadataRoundTrip тестирует полный цикл упаковки-распаковки метаданных
func TestPackUnpackMetadataRoundTrip(t *testing.T) {
	codec := codec.NewCodec()

	// Тестируем UserMetadata
	metadata := &vault.UserMetadata{
		Version:      "1.0",
		Settings:     map[string]string{"theme": "dark", "language": "en"},
		PasswordIDs:  []string{"id1", "id2", "id3"},
		UpdatedAt:    time.Now(),
		TotalEntries: 3,
	}

	masterPassword := "test-master-password-123"

	// Упаковываем
	packedData, err := codec.PackMetadata(metadata, masterPassword)
	if err != nil {
		t.Fatalf("PackMetadata failed: %v", err)
	}

	// Распаковываем
	unpackedMetadata, err := codec.UnpackMetadata(packedData, masterPassword)
	if err != nil {
		t.Fatalf("UnpackMetadata failed: %v", err)
	}

	// Проверяем, что данные идентичны (кроме PasswordIDs, которые не сохраняются)
	if unpackedMetadata.Version != metadata.Version {
		t.Error("Version should be identical after round trip")
	}
	if unpackedMetadata.TotalEntries != metadata.TotalEntries {
		t.Error("TotalEntries should be identical after round trip")
	}

	// Проверяем настройки
	for key, value := range metadata.Settings {
		if unpackedMetadata.Settings[key] != value {
			t.Errorf("Setting %s should be identical after round trip", key)
		}
	}
}
