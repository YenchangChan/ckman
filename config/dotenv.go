package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnvFile loads conf/.ckman.env beside the selected ckman config file.
// An absent file is optional; malformed or unreadable files are errors.
// godotenv.Load preserves values already supplied by the process environment.
func LoadEnvFile(configFile string) error {
	envFile := filepath.Join(filepath.Dir(configFile), ".ckman.env")
	if err := godotenv.Load(envFile); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("load env file %s: %w", envFile, err)
	}
	return nil
}
