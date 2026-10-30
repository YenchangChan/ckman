package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadEnvFileBesideConfig(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "ckman.hjson")
	envFile := filepath.Join(dir, ".ckman.env")
	const key = "CKMAN_DOTENV_TEST_SECRET"
	old, existed := os.LookupEnv(key)
	require.NoError(t, os.Unsetenv(key))
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	})
	require.NoError(t, LoadEnvFile(configFile)) // optional when absent
	require.NoError(t, os.WriteFile(envFile, []byte(key+"=from-file\n"), 0600))
	require.NoError(t, LoadEnvFile(configFile))
	require.Equal(t, "from-file", os.Getenv(key))
	require.NoError(t, os.Setenv(key, "from-process"))
	require.NoError(t, LoadEnvFile(configFile))
	require.Equal(t, "from-process", os.Getenv(key))
}

func TestParseConfigFileLoadsDotEnv(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "ckman.hjson")
	const key = "CKMAN_DOTENV_PARSE_TEST"
	old, existed := os.LookupEnv(key)
	require.NoError(t, os.Unsetenv(key))
	previousConfig := GlobalConfig
	t.Cleanup(func() {
		GlobalConfig = previousConfig
		if existed {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	})
	require.NoError(t, os.WriteFile(configFile, []byte("{}"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ckman.env"), []byte(key+"=loaded\n"), 0600))
	require.NoError(t, ParseConfigFile(configFile, ""))
	require.Equal(t, "loaded", os.Getenv(key))
}
