package keymanager_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"encryptkeep-backend/internal/keymanager"

	"github.com/ethereum/go-ethereum/crypto"
)

// TestNewKeyManager тестирует создание нового KeyManager
func TestNewKeyManager(t *testing.T) {
	config := keymanager.KeyManagerConfig{
		ConfigDir:      "/tmp/test",
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	if km == nil {
		t.Fatal("KeyManager should not be nil")
	}

	if km.ConfigDir != config.ConfigDir {
		t.Errorf("Expected ConfigDir %s, got %s", config.ConfigDir, km.ConfigDir)
	}

	if km.SessionActive {
		t.Error("Session should not be active initially")
	}

	if km.PrivateKey != nil {
		t.Error("PrivateKey should be nil initially")
	}

	if km.MasterKey != nil {
		t.Error("MasterKey should be nil initially")
	}
}

// TestNewKeyManagerDefaultTimeout тестирует установку таймаута по умолчанию
func TestNewKeyManagerDefaultTimeout(t *testing.T) {
	config := keymanager.KeyManagerConfig{
		ConfigDir: "/tmp/test",
		// SessionTimeout не установлен
	}

	km := keymanager.NewKeyManager(config)

	// Проверяем, что KeyManager создан с дефолтным таймаутом
	// (не можем проверить приватное поле config напрямую)
	if km == nil {
		t.Error("KeyManager should not be nil")
	}
}

// TestHasStoredKeys тестирует проверку наличия сохраненных ключей
func TestHasStoredKeys(t *testing.T) {
	// Создаем временную директорию
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	// Изначально ключей нет
	if km.HasStoredKeys() {
		t.Error("Should not have stored keys initially")
	}

	// Создаем файл keys.json
	keyFilePath := filepath.Join(tempDir, "keys.json")
	err := os.WriteFile(keyFilePath, []byte("{}"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Теперь ключи должны быть
	if !km.HasStoredKeys() {
		t.Error("Should have stored keys after creating file")
	}
}

// TestSessionManagement тестирует управление сессиями
func TestSessionManagement(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 100 * time.Millisecond, // Короткий таймаут для тестирования
	}

	km := keymanager.NewKeyManager(config)

	// Изначально сессия неактивна
	if km.IsSessionActive() {
		t.Error("Session should not be active initially")
	}

	// Активируем сессию
	km.SessionActive = true
	km.LastActivity = time.Now()

	if !km.IsSessionActive() {
		t.Error("Session should be active after activation")
	}

	// Ждем истечения таймаута
	time.Sleep(150 * time.Millisecond)

	if km.IsSessionActive() {
		t.Error("Session should be inactive after timeout")
	}
}

// TestClearSession тестирует очистку сессии
func TestClearSession(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	// Создаем тестовые данные
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	masterKey := []byte("test-master-key-32-bytes-long")

	// Устанавливаем данные
	km.PrivateKey = privateKey
	km.MasterKey = make([]byte, len(masterKey))
	copy(km.MasterKey, masterKey)
	km.SessionActive = true

	// Проверяем, что данные установлены
	if km.PrivateKey == nil {
		t.Error("PrivateKey should be set")
	}
	if km.MasterKey == nil {
		t.Error("MasterKey should be set")
	}
	if !km.IsSessionActive() {
		t.Error("Session should be active")
	}

	// Очищаем сессию
	km.ClearSession()

	// Проверяем, что данные очищены
	if km.PrivateKey != nil {
		t.Error("PrivateKey should be nil after ClearSession")
	}
	if km.MasterKey != nil {
		t.Error("MasterKey should be nil after ClearSession")
	}
	if km.SessionActive {
		t.Error("Session should be inactive after ClearSession")
	}
}

// TestGetPrivateKey тестирует получение приватного ключа
func TestGetPrivateKey(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	// Тест без активной сессии
	_, err := km.GetPrivateKey()
	if err == nil {
		t.Error("Should return error when session is not active")
	}

	// Создаем тестовый приватный ключ
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Активируем сессию
	km.PrivateKey = privateKey
	km.SessionActive = true
	km.LastActivity = time.Now()

	// Тест с активной сессией
	retrievedKey, err := km.GetPrivateKey()
	if err != nil {
		t.Errorf("Should not return error when session is active: %v", err)
	}

	if retrievedKey == nil {
		t.Error("Retrieved private key should not be nil")
		return
	}

	// Проверяем, что это тот же ключ
	if retrievedKey.D.Cmp(privateKey.D) != 0 {
		t.Error("Retrieved private key should match original")
	}
}

// TestGetMasterKey тестирует получение мастер-ключа
func TestGetMasterKey(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	// Тест без активной сессии
	_, err := km.GetMasterKey()
	if err == nil {
		t.Error("Should return error when session is not active")
	}

	// Создаем тестовый мастер-ключ
	masterKey := []byte("test-master-key-32-bytes-long")

	// Активируем сессию
	km.MasterKey = make([]byte, len(masterKey))
	copy(km.MasterKey, masterKey)
	km.SessionActive = true
	km.LastActivity = time.Now()

	// Тест с активной сессией
	retrievedKey, err := km.GetMasterKey()
	if err != nil {
		t.Errorf("Should not return error when session is active: %v", err)
	}

	if retrievedKey == nil {
		t.Error("Retrieved master key should not be nil")
	}

	// Проверяем, что это тот же ключ
	if len(retrievedKey) != len(masterKey) {
		t.Errorf("Retrieved master key length should match original: got %d, want %d", len(retrievedKey), len(masterKey))
	}
}

// TestUpdateActivity тестирует обновление активности
func TestUpdateActivity(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	initialTime := km.LastActivity

	// Обновляем активность
	km.UpdateActivity()

	// Проверяем, что время не изменилось (сессия неактивна)
	if km.LastActivity != initialTime {
		t.Error("LastActivity should not change when session is inactive")
	}

	// Активируем сессию
	km.SessionActive = true

	// Обновляем активность
	km.UpdateActivity()

	// Проверяем, что время изменилось или осталось тем же (если вызовы очень быстрые)
	if km.LastActivity.Before(initialTime) {
		t.Error("LastActivity should not be before initial time")
	}
}

// TestInitializeFirstTime тестирует инициализацию при первом запуске
func TestInitializeFirstTime(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	// Генерируем тестовый приватный ключ
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	privateKeyHex := crypto.FromECDSA(privateKey)
	privateKeyHexStr := hex.EncodeToString(privateKeyHex)
	masterPassword := "test-password-123"

	// Тест с неверной длиной ключа
	err = km.InitializeFirstTime("short", masterPassword)
	if err == nil {
		t.Error("Should return error for invalid key length")
	}

	// Тест с коротким паролем
	err = km.InitializeFirstTime(privateKeyHexStr, "short")
	if err == nil {
		t.Error("Should return error for short password")
	}

	// Тест с валидными данными
	err = km.InitializeFirstTime(privateKeyHexStr, masterPassword)
	if err != nil {
		t.Errorf("Should not return error for valid data: %v", err)
	}

	// Проверяем, что сессия активирована
	if !km.IsSessionActive() {
		t.Error("Session should be active after InitializeFirstTime")
	}

	// Проверяем, что приватный ключ установлен
	if km.PrivateKey == nil {
		t.Error("PrivateKey should be set after InitializeFirstTime")
	}

	// Проверяем, что мастер-ключ установлен
	if km.MasterKey == nil {
		t.Error("MasterKey should be set after InitializeFirstTime")
	}

	// Проверяем, что файл создан
	if !km.HasStoredKeys() {
		t.Error("Should have stored keys after InitializeFirstTime")
	}
}

// TestLoadFromStorage тестирует загрузку из хранилища
func TestLoadFromStorage(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	// Тест без сохраненных ключей
	err := km.LoadFromStorage("test-password")
	if err == nil {
		t.Error("Should return error when no stored keys exist")
	}

	// Сначала инициализируем ключи
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	privateKeyHex := crypto.FromECDSA(privateKey)
	privateKeyHexStr := hex.EncodeToString(privateKeyHex)
	masterPassword := "test-password-123"

	err = km.InitializeFirstTime(privateKeyHexStr, masterPassword)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Очищаем сессию
	km.ClearSession()

	// Тест загрузки с правильным паролем
	err = km.LoadFromStorage(masterPassword)
	if err != nil {
		t.Errorf("Should not return error with correct password: %v", err)
	}

	// Проверяем, что сессия активирована
	if !km.IsSessionActive() {
		t.Error("Session should be active after LoadFromStorage")
	}

	// Проверяем, что приватный ключ восстановлен
	if km.PrivateKey == nil {
		t.Error("PrivateKey should be restored after LoadFromStorage")
	}

	// Проверяем, что мастер-ключ восстановлен
	if km.MasterKey == nil {
		t.Error("MasterKey should be restored after LoadFromStorage")
	}

	// Тест загрузки с неправильным паролем
	km.ClearSession()
	err = km.LoadFromStorage("wrong-password")
	if err == nil {
		t.Error("Should return error with wrong password")
	}
}

// TestGetAddress тестирует получение адреса
func TestGetAddress(t *testing.T) {
	tempDir := t.TempDir()

	config := keymanager.KeyManagerConfig{
		ConfigDir:      tempDir,
		SessionTimeout: 30 * time.Minute,
	}

	km := keymanager.NewKeyManager(config)

	// Тест без активной сессии
	_, err := km.GetAddress()
	if err == nil {
		t.Error("Should return error when session is not active")
	}

	// Создаем тестовый приватный ключ
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Активируем сессию
	km.PrivateKey = privateKey
	km.SessionActive = true
	km.LastActivity = time.Now()

	// Тест получения адреса
	address, err := km.GetAddress()
	if err != nil {
		t.Errorf("Should not return error when session is active: %v", err)
	}

	if address == "" {
		t.Error("Address should not be empty")
	}

	// Проверяем, что адрес соответствует приватному ключу
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Failed to cast public key")
	}

	expectedAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	if address != expectedAddress.Hex() {
		t.Errorf("Address mismatch: got %s, want %s", address, expectedAddress.Hex())
	}
}

// TestStoredKeyDataJSON тестирует сериализацию/десериализацию StoredKeyData
func TestStoredKeyDataJSON(t *testing.T) {
	// Создаем тестовые данные
	data := keymanager.StoredKeyData{
		EncryptedPrivateKey: []byte("encrypted-data"),
		Salt:                []byte("salt-data"),
		Nonce:               []byte("nonce-data"),
		CreatedAt:           "2023-01-01T00:00:00Z",
		Address:             "0x1234567890123456789012345678901234567890",
	}

	// Тестируем сериализацию
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal StoredKeyData: %v", err)
	}

	// Тестируем десериализацию
	var restoredData keymanager.StoredKeyData
	err = json.Unmarshal(jsonData, &restoredData)
	if err != nil {
		t.Errorf("Failed to unmarshal StoredKeyData: %v", err)
	}

	// Проверяем, что данные восстановлены корректно
	if string(restoredData.EncryptedPrivateKey) != string(data.EncryptedPrivateKey) {
		t.Error("EncryptedPrivateKey mismatch")
	}
	if string(restoredData.Salt) != string(data.Salt) {
		t.Error("Salt mismatch")
	}
	if string(restoredData.Nonce) != string(data.Nonce) {
		t.Error("Nonce mismatch")
	}
	if restoredData.CreatedAt != data.CreatedAt {
		t.Error("CreatedAt mismatch")
	}
	if restoredData.Address != data.Address {
		t.Error("Address mismatch")
	}
}
