package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersionedSecretsAndLegacyAES(t *testing.T) {
	t.Setenv(SM4KeyEnv, "0123456789abcdef")
	aes, err := EncryptSecret("legacy", GosyptAlgorithm)
	require.NoError(t, err)
	require.NotContains(t, aes, SM4SecretPrefix)
	sm4, err := EncryptSecret("new-secret", GosyptSM4Algorithm)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(sm4, SM4SecretPrefix))
	for _, tc := range []struct{ cipher, plain string }{{aes, "legacy"}, {sm4, "new-secret"}} {
		got, err := DecryptSecret(tc.cipher)
		require.NoError(t, err)
		require.Equal(t, tc.plain, got)
	}
	g := &Gosypt{prefix: GosyptPrefixDefault, suffix: GosyptSuffxiDefault}
	require.NoError(t, g.SetAlgorithm(GosyptSM4Algorithm))
	for _, tc := range []struct{ cipher, plain string }{{aes, "legacy"}, {sm4, "new-secret"}} {
		cfg := struct{ Password string }{Password: "ENC(" + tc.cipher + ")"}
		require.NoError(t, g.Unmarshal(&cfg))
		require.Equal(t, tc.plain, cfg.Password)
	}
	t.Setenv(SM4KeyEnv, "fedcba9876543210")
	_, err = DecryptSecret(sm4)
	require.Error(t, err)
	_, err = DecryptSecret("sm4:gcm:v2:anything")
	require.Error(t, err)
	_, err = DecryptSecret("not-ciphertext")
	require.Error(t, err)
}
