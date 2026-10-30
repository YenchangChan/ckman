package common

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const SM4SecretPrefix = "sm4:gcm:v1:"

// EncryptSecret uses the configured algorithm for new secrets. The prefix makes
// SM4 ciphertext readable even after the write algorithm is changed.
func EncryptSecret(plain, algorithm string) (string, error) {
	if plain == "" {
		return "", nil
	}
	switch algorithm {
	case "", GosyptAlgorithm:
		return AesEncryptECB(plain), nil
	case GosyptSM4Algorithm:
		data, err := SM4Encrypt(plain)
		if err != nil {
			return "", err
		}
		return SM4SecretPrefix + data, nil
	default:
		return "", fmt.Errorf("unsupported encryption algorithm %q", algorithm)
	}
}

// DecryptSecret reads both versioned SM4 and legacy unmarked AES ciphertext.
func DecryptSecret(data string) (string, error) {
	if data == "" {
		return "", nil
	}
	if strings.HasPrefix(data, SM4SecretPrefix) {
		payload := strings.TrimPrefix(data, SM4SecretPrefix)
		if payload == "" {
			return "", fmt.Errorf("empty SM4 ciphertext")
		}
		return SM4Decrypt(payload)
	}
	if strings.HasPrefix(data, "sm4:") {
		return "", fmt.Errorf("unsupported SM4 ciphertext format")
	}
	if len(data)%32 != 0 {
		return "", fmt.Errorf("invalid legacy AES ciphertext length")
	}
	if _, err := hex.DecodeString(data); err != nil {
		return "", fmt.Errorf("invalid legacy AES ciphertext: %w", err)
	}
	plain := AesDecryptECB(data)
	if AesEncryptECB(plain) != strings.ToUpper(data) {
		return "", fmt.Errorf("legacy AES ciphertext verification failed")
	}
	return plain, nil
}

func looksLikeLegacyAES(data string) bool {
	if len(data) == 0 || len(data)%32 != 0 {
		return false
	}
	_, err := hex.DecodeString(data)
	return err == nil
}
