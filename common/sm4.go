package common

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/tjfoc/gmsm/sm4"
)

const (
	SM4KeyEnv     = "CKMAN_SM4_KEY"
	SM4KeyFileEnv = "CKMAN_SM4_KEY_FILE"
)

// GenerateSM4KeyHex returns a cryptographically random 16-byte key in hex form.
func GenerateSM4KeyHex() (string, error) {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func sm4Key() ([]byte, error) {
	value := os.Getenv(SM4KeyEnv)
	if value == "" {
		if path := os.Getenv(SM4KeyFileEnv); path != "" {
			key, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("read %s file %s: %w", SM4KeyFileEnv, path, err)
			}
			if len(key) != 16 {
				return nil, fmt.Errorf("%s file must contain exactly 16 raw bytes", SM4KeyFileEnv)
			}
			return key, nil
		}
		return nil, fmt.Errorf("set %s (16 raw bytes, 32 hex characters, hex:<32 hex>, or base64:<16-byte key>) or %s (path to 16-byte raw key file)", SM4KeyEnv, SM4KeyFileEnv)
	}
	// Preserve existing 16-byte raw keys, even if their text starts with a
	// representation prefix such as "hex:".
	if len(value) == 16 {
		return []byte(value), nil
	}
	if strings.HasPrefix(value, "base64:") {
		key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, "base64:"))
		if err != nil || len(key) != 16 {
			return nil, fmt.Errorf("%s base64 value must decode to exactly 16 bytes", SM4KeyEnv)
		}
		return key, nil
	}
	if strings.HasPrefix(value, "hex:") {
		value = strings.TrimPrefix(value, "hex:")
		if len(value) != 32 {
			return nil, fmt.Errorf("%s hex value must contain exactly 32 hexadecimal characters", SM4KeyEnv)
		}
	}
	if len(value) == 32 {
		key, err := hex.DecodeString(value)
		if err != nil {
			return nil, fmt.Errorf("%s must be valid hexadecimal: %w", SM4KeyEnv, err)
		}
		return key, nil
	}
	return nil, fmt.Errorf("%s must be 16 raw bytes or 32 hexadecimal characters; use ckmanctl keygen sm4 for a new key, or base64:<value> for an existing binary key", SM4KeyEnv)
}

func SM4Encrypt(data string) (string, error) {
	if data == "" {
		return "", nil
	}
	key, err := sm4Key()
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(data), nil)), nil
}

func SM4Decrypt(data string) (string, error) {
	if data == "" {
		return "", nil
	}
	key, err := sm4Key()
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	if len(decoded) < gcm.NonceSize() {
		return "", fmt.Errorf("SM4 ciphertext is too short")
	}
	plain, err := gcm.Open(nil, decoded[:gcm.NonceSize()], decoded[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
