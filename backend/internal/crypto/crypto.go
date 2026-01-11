package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const DefaultSaltLen = 16
const GCMNonceLen = 12
const AESKeyLen = 32

type Argon2Config struct {
	Time      uint32
	Memory    uint32
	Threads   uint8
	KeyLength uint32
}

type DerivedKey struct {
	Key  []byte // длина = KeyLength
	Salt []byte // 16 байт по умолчанию
}

type Sealed struct {
	Salt       []byte
	Nonce      []byte
	Ciphertext []byte
}

func DeriveKey(masterPassword string, salt []byte, cfg Argon2Config) (DerivedKey, error) {
	if masterPassword == "" {
		return DerivedKey{}, errors.New("master password must not be empty")
	}

	// Basic config validation
	if cfg.Time == 0 {
		return DerivedKey{}, errors.New("argon2 time parameter must be > 0")
	}
	if cfg.Memory == 0 {
		return DerivedKey{}, errors.New("argon2 memory parameter must be > 0")
	}
	if cfg.Threads == 0 {
		return DerivedKey{}, errors.New("argon2 threads parameter must be > 0")
	}
	if cfg.KeyLength == 0 {
		cfg.KeyLength = AESKeyLen
	}

	// Ensure salt is present
	var err error
	if len(salt) == 0 {
		salt, err = GenerateSalt(DefaultSaltLen)
		if err != nil {
			return DerivedKey{}, err
		}
	}

	key := argon2.IDKey([]byte(masterPassword), salt, cfg.Time, cfg.Memory, cfg.Threads, cfg.KeyLength)
	if len(key) == 0 {
		return DerivedKey{}, errors.New("derived key is empty")
	}

	return DerivedKey{Key: key, Salt: salt}, nil
}

func Encrypt(key, plaintext []byte) (ciphertext, nonce []byte, err error) {
	if len(key) != AESKeyLen {
		return nil, nil, fmt.Errorf("invalid key length: got %d, want %d", len(key), AESKeyLen)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("new gcm: %w", err)
	}

	nonce, err = GenerateNonce()
	if err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func Decrypt(key, ciphertext, nonce []byte) ([]byte, error) {
	if len(key) != AESKeyLen {
		return nil, fmt.Errorf("invalid key length: got %d, want %d", len(key), AESKeyLen)
	}
	if len(nonce) != GCMNonceLen {
		return nil, fmt.Errorf("invalid nonce length: got %d, want %d", len(nonce), GCMNonceLen)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm open: %w", err)
	}
	return plaintext, nil
}

func GenerateSalt(n int) ([]byte, error) {
	if n <= 0 {
		return nil, errors.New("salt length must be positive")
	}

	salt := make([]byte, n)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	return salt, nil
}

func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, GCMNonceLen)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return nonce, nil
}

func Seal(masterPassword string, cfg Argon2Config, plaintext []byte) (Sealed, error) {
	dk, err := DeriveKey(masterPassword, nil, cfg)
	if err != nil {
		return Sealed{}, err
	}

	ciphertext, nonce, err := Encrypt(dk.Key, plaintext)
	if err != nil {
		return Sealed{}, err
	}

	return Sealed{
		Salt:       dk.Salt,
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}, nil
}

func Open(masterPassword string, cfg Argon2Config, sealed Sealed) ([]byte, error) {
	if len(sealed.Salt) == 0 {
		return nil, errors.New("sealed salt must not be empty")
	}
	dk, err := DeriveKey(masterPassword, sealed.Salt, cfg)
	if err != nil {
		return nil, err
	}

	plaintext, err := Decrypt(dk.Key, sealed.Ciphertext, sealed.Nonce)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
