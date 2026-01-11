package crypto_test

import (
	"testing"

	"encryptkeep-backend/internal/crypto"
)

// TestArgon2Config тестирует структуру Argon2Config
func TestArgon2Config(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	if config.Time == 0 {
		t.Error("Time should not be zero")
	}
	if config.Memory == 0 {
		t.Error("Memory should not be zero")
	}
	if config.Threads == 0 {
		t.Error("Threads should not be zero")
	}
	if config.KeyLength == 0 {
		t.Error("KeyLength should not be zero")
	}
}

// TestDeriveKey тестирует деривацию ключа
func TestDeriveKey(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	masterPassword := "test-password-123"

	// Тест с пустым паролем
	_, err := crypto.DeriveKey("", nil, config)
	if err == nil {
		t.Error("Should return error for empty password")
	}

	// Тест с валидным паролем
	dk, err := crypto.DeriveKey(masterPassword, nil, config)
	if err != nil {
		t.Errorf("Should not return error for valid password: %v", err)
	}

	if dk.Key == nil {
		t.Error("Derived key should not be nil")
	}

	if len(dk.Key) != int(config.KeyLength) {
		t.Errorf("Key length mismatch: got %d, want %d", len(dk.Key), config.KeyLength)
	}

	if dk.Salt == nil {
		t.Error("Salt should not be nil")
	}

	if len(dk.Salt) == 0 {
		t.Error("Salt should not be empty")
	}
}

// TestDeriveKeyWithSalt тестирует деривацию ключа с заданной солью
func TestDeriveKeyWithSalt(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	masterPassword := "test-password-123"
	customSalt := []byte("custom-salt-16")

	dk, err := crypto.DeriveKey(masterPassword, customSalt, config)
	if err != nil {
		t.Errorf("Should not return error: %v", err)
	}

	if dk.Key == nil {
		t.Error("Derived key should not be nil")
	}

	if len(dk.Key) != int(config.KeyLength) {
		t.Errorf("Key length mismatch: got %d, want %d", len(dk.Key), config.KeyLength)
	}

	// Проверяем, что используется заданная соль
	if len(dk.Salt) != len(customSalt) {
		t.Errorf("Salt length mismatch: got %d, want %d", len(dk.Salt), len(customSalt))
	}
}

// TestDeriveKeyInvalidConfig тестирует деривацию с невалидным конфигом
func TestDeriveKeyInvalidConfig(t *testing.T) {
	masterPassword := "test-password-123"

	// Тест с нулевым Time
	config := crypto.Argon2Config{
		Time:      0, // Невалидно
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	_, err := crypto.DeriveKey(masterPassword, nil, config)
	if err == nil {
		t.Error("Should return error for zero Time")
	}

	// Тест с нулевой Memory
	config = crypto.Argon2Config{
		Time:      3,
		Memory:    0, // Невалидно
		Threads:   4,
		KeyLength: 32,
	}

	_, err = crypto.DeriveKey(masterPassword, nil, config)
	if err == nil {
		t.Error("Should return error for zero Memory")
	}

	// Тест с нулевыми Threads
	config = crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   0, // Невалидно
		KeyLength: 32,
	}

	_, err = crypto.DeriveKey(masterPassword, nil, config)
	if err == nil {
		t.Error("Should return error for zero Threads")
	}
}

// TestEncryptDecrypt тестирует шифрование и расшифровку
func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, crypto.AESKeyLen)
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := []byte("Hello, World! This is a test message.")

	// Тест шифрования
	ciphertext, nonce, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Errorf("Encryption failed: %v", err)
	}

	if ciphertext == nil {
		t.Error("Ciphertext should not be nil")
	}

	if nonce == nil {
		t.Error("Nonce should not be nil")
	}

	if len(nonce) != crypto.GCMNonceLen {
		t.Errorf("Nonce length mismatch: got %d, want %d", len(nonce), crypto.GCMNonceLen)
	}

	// Тест расшифровки
	decrypted, err := crypto.Decrypt(key, ciphertext, nonce)
	if err != nil {
		t.Errorf("Decryption failed: %v", err)
	}

	if decrypted == nil {
		t.Error("Decrypted text should not be nil")
	}

	// Проверяем, что расшифрованный текст совпадает с оригиналом
	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted text mismatch: got %s, want %s", string(decrypted), string(plaintext))
	}
}

// TestEncryptInvalidKey тестирует шифрование с невалидным ключом
func TestEncryptInvalidKey(t *testing.T) {
	invalidKey := []byte("short-key") // Неправильная длина
	plaintext := []byte("test message")

	_, _, err := crypto.Encrypt(invalidKey, plaintext)
	if err == nil {
		t.Error("Should return error for invalid key length")
	}
}

// TestDecryptInvalidKey тестирует расшифровку с невалидным ключом
func TestDecryptInvalidKey(t *testing.T) {
	invalidKey := []byte("short-key") // Неправильная длина
	ciphertext := []byte("encrypted data")
	nonce := make([]byte, crypto.GCMNonceLen)

	_, err := crypto.Decrypt(invalidKey, ciphertext, nonce)
	if err == nil {
		t.Error("Should return error for invalid key length")
	}
}

// TestDecryptInvalidNonce тестирует расшифровку с невалидным nonce
func TestDecryptInvalidNonce(t *testing.T) {
	key := make([]byte, crypto.AESKeyLen)
	ciphertext := []byte("encrypted data")
	invalidNonce := []byte("short-nonce") // Неправильная длина

	_, err := crypto.Decrypt(key, ciphertext, invalidNonce)
	if err == nil {
		t.Error("Should return error for invalid nonce length")
	}
}

// TestGenerateSalt тестирует генерацию соли
func TestGenerateSalt(t *testing.T) {
	// Тест с положительной длиной
	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		t.Errorf("Should not return error: %v", err)
	}

	if salt == nil {
		t.Error("Salt should not be nil")
	}

	if len(salt) != 16 {
		t.Errorf("Salt length mismatch: got %d, want 16", len(salt))
	}

	// Тест с нулевой длиной
	_, err = crypto.GenerateSalt(0)
	if err == nil {
		t.Error("Should return error for zero length")
	}

	// Тест с отрицательной длиной
	_, err = crypto.GenerateSalt(-1)
	if err == nil {
		t.Error("Should return error for negative length")
	}
}

// TestGenerateNonce тестирует генерацию nonce
func TestGenerateNonce(t *testing.T) {
	nonce, err := crypto.GenerateNonce()
	if err != nil {
		t.Errorf("Should not return error: %v", err)
	}

	if nonce == nil {
		t.Error("Nonce should not be nil")
	}

	if len(nonce) != crypto.GCMNonceLen {
		t.Errorf("Nonce length mismatch: got %d, want %d", len(nonce), crypto.GCMNonceLen)
	}
}

// TestSealOpen тестирует функции Seal и Open
func TestSealOpen(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	masterPassword := "test-password-123"
	plaintext := []byte("This is a secret message that needs to be encrypted.")

	// Тест Seal
	sealed, err := crypto.Seal(masterPassword, config, plaintext)
	if err != nil {
		t.Errorf("Seal failed: %v", err)
	}

	if sealed.Salt == nil {
		t.Error("Sealed Salt should not be nil")
	}

	if sealed.Nonce == nil {
		t.Error("Sealed Nonce should not be nil")
	}

	if sealed.Ciphertext == nil {
		t.Error("Sealed Ciphertext should not be nil")
	}

	// Тест Open
	opened, err := crypto.Open(masterPassword, config, sealed)
	if err != nil {
		t.Errorf("Open failed: %v", err)
	}

	if opened == nil {
		t.Error("Opened data should not be nil")
	}

	// Проверяем, что открытые данные совпадают с оригиналом
	if string(opened) != string(plaintext) {
		t.Errorf("Opened data mismatch: got %s, want %s", string(opened), string(plaintext))
	}
}

// TestSealOpenWrongPassword тестирует Open с неправильным паролем
func TestSealOpenWrongPassword(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	masterPassword := "test-password-123"
	wrongPassword := "wrong-password-456"
	plaintext := []byte("This is a secret message.")

	// Seal с правильным паролем
	sealed, err := crypto.Seal(masterPassword, config, plaintext)
	if err != nil {
		t.Errorf("Seal failed: %v", err)
	}

	// Open с неправильным паролем
	_, err = crypto.Open(wrongPassword, config, sealed)
	if err == nil {
		t.Error("Should return error for wrong password")
	}
}

// TestSealOpenEmptyPassword тестирует Seal с пустым паролем
func TestSealOpenEmptyPassword(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	plaintext := []byte("test message")

	_, err := crypto.Seal("", config, plaintext)
	if err == nil {
		t.Error("Should return error for empty password")
	}
}

// TestSealOpenEmptySealed тестирует Open с пустой структурой Sealed
func TestSealOpenEmptySealed(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	emptySealed := crypto.Sealed{
		Salt:       nil, // Пустая соль
		Nonce:      []byte("nonce"),
		Ciphertext: []byte("ciphertext"),
	}

	_, err := crypto.Open("password", config, emptySealed)
	if err == nil {
		t.Error("Should return error for empty salt")
	}
}

// TestDerivedKeyConsistency тестирует консистентность деривации ключа
func TestDerivedKeyConsistency(t *testing.T) {
	config := crypto.Argon2Config{
		Time:      3,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	masterPassword := "test-password-123"
	salt := []byte("test-salt-16-by")

	// Деривируем ключ дважды с одинаковыми параметрами
	dk1, err1 := crypto.DeriveKey(masterPassword, salt, config)
	dk2, err2 := crypto.DeriveKey(masterPassword, salt, config)

	if err1 != nil {
		t.Errorf("First derivation failed: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Second derivation failed: %v", err2)
	}

	// Ключи должны быть одинаковыми
	if len(dk1.Key) != len(dk2.Key) {
		t.Error("Key lengths should be equal")
	}

	for i := range dk1.Key {
		if dk1.Key[i] != dk2.Key[i] {
			t.Error("Derived keys should be identical")
			break
		}
	}

	// Соли должны быть одинаковыми
	if len(dk1.Salt) != len(dk2.Salt) {
		t.Error("Salt lengths should be equal")
	}

	for i := range dk1.Salt {
		if dk1.Salt[i] != dk2.Salt[i] {
			t.Error("Salts should be identical")
			break
		}
	}
}
