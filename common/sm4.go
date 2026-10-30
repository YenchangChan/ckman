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

// GenerateSM4KeyBase64 returns a cryptographically random 16-byte key in Base64 form.
func GenerateSM4KeyBase64() (string, error) {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
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
		return nil, fmt.Errorf("set %s to the Base64 encoding of a 16-byte SM4 key, or set %s to a 16-byte raw key file", SM4KeyEnv, SM4KeyFileEnv)
	}
	// Preserve existing 16-byte raw keys, even if their text starts with a
	// representation prefix such as "hex:".
	if len(value) == 16 {
		return []byte(value), nil
	}
	if strings.HasPrefix(value, "hex:") {
		value = strings.TrimPrefix(value, "hex:")
		if len(value) != 32 {
			return nil, fmt.Errorf("%s hex value must contain exactly 32 hexadecimal characters", SM4KeyEnv)
		}
		key, err := hex.DecodeString(value)
		if err != nil {
			return nil, fmt.Errorf("%s hex value is invalid: %w", SM4KeyEnv, err)
		}
		return key, nil
	}
	if len(value) == 32 {
		if key, err := hex.DecodeString(value); err == nil {
			return key, nil
		}
	}
	value = strings.TrimPrefix(value, "base64:")
	for _, encoding := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		if key, err := encoding.DecodeString(value); err == nil && len(key) == 16 {
			return key, nil
		}
	}
	return nil, fmt.Errorf("%s must be Base64 text that decodes to exactly 16 bytes (or a legacy 16-byte raw/32-character hex key)", SM4KeyEnv)
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
