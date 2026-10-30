package common

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateSM4KeyBase64(t *testing.T) {
	first, err := GenerateSM4KeyBase64()
	require.NoError(t, err)
	second, err := GenerateSM4KeyBase64()
	require.NoError(t, err)
	require.NotEqual(t, first, second)
	require.Len(t, first, 24)
	key, err := base64.StdEncoding.DecodeString(first)
	require.NoError(t, err)
	require.Len(t, key, 16)
}

func TestSM4ExistingKeyRepresentations(t *testing.T) {
	key := []byte{'a', 'b', 'c', 'd', 0, 1, 2, 3, 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l'}
	path := filepath.Join(t.TempDir(), "sm4.key")
	require.NoError(t, os.WriteFile(path, key, 0600))
	for _, value := range []string{
		base64.StdEncoding.EncodeToString(key),
		base64.RawStdEncoding.EncodeToString(key),
		base64.URLEncoding.EncodeToString(key),
		"base64:" + base64.StdEncoding.EncodeToString(key),
		"hex:" + hex.EncodeToString(key),
		hex.EncodeToString(key),
	} {
		t.Setenv(SM4KeyEnv, value)
		got, err := sm4Key()
		require.NoError(t, err)
		require.Equal(t, key, got)
	}
	t.Setenv(SM4KeyEnv, "")
	t.Setenv(SM4KeyFileEnv, path)
	got, err := sm4Key()
	require.NoError(t, err)
	require.Equal(t, key, got)
	ciphertext, err := SM4Encrypt("existing-key")
	require.NoError(t, err)
	t.Setenv(SM4KeyEnv, base64.StdEncoding.EncodeToString(key))
	plain, err := SM4Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, "existing-key", plain)
	t.Setenv(SM4KeyEnv, "base64:invalid")
	_, err = sm4Key()
	require.Error(t, err)
	t.Setenv(SM4KeyEnv, "hex:123456789012")
	got, err = sm4Key()
	require.NoError(t, err)
	require.Equal(t, []byte("hex:123456789012"), got)
}

func TestSM4OriginalKeyFromHexEnv(t *testing.T) {
	t.Setenv(SM4KeyEnv, "6368656e79616e6368616e6702000205")
	key, err := sm4Key()
	require.NoError(t, err)
	require.Equal(t, []byte{'c', 'h', 'e', 'n', 'y', 'a', 'n', 'c', 'h', 'a', 'n', 'g', 0x02, 0x00, 0x02, 0x05}, key)
	ciphertext, err := SM4Encrypt("test")
	require.NoError(t, err)
	plain, err := SM4Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, "test", plain)
}

func TestSM4ConfigSecrets(t *testing.T) {
	t.Setenv(SM4KeyEnv, "1234567890abcdef")
	first, err := SM4Encrypt("secret-value")
	require.NoError(t, err)
	second, err := SM4Encrypt("secret-value")
	require.NoError(t, err)
	require.NotEqual(t, first, second)
	got, err := SM4Decrypt(first)
	require.NoError(t, err)
	require.Equal(t, "secret-value", got)

	g := &Gosypt{prefix: GosyptPrefixDefault, suffix: GosyptSuffxiDefault}
	require.NoError(t, g.SetAlgorithm(GosyptSM4Algorithm))
	cfg := struct{ Password string }{Password: "ENC(" + first + ")"}
	require.NoError(t, g.Unmarshal(&cfg))
	require.Equal(t, "secret-value", cfg.Password)

	nested := map[string]interface{}{"database": map[string]interface{}{"password": "ENC(" + first + ")"}}
	require.NoError(t, g.Unmarshal(&nested))
	require.Equal(t, "secret-value", nested["database"].(map[string]interface{})["password"])

	cfg.Password = "ENC(corrupted)"
	require.Error(t, g.Unmarshal(&cfg))
	t.Setenv(SM4KeyEnv, "wrong-length")
	require.Error(t, g.SetAlgorithm(GosyptSM4Algorithm))
	require.Error(t, g.Unmarshal(&cfg))
}
