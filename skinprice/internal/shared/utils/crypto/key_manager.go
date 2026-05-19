package crypto

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

const tokenKeyFileName = "token_encryption.key"

func LoadOrCreateTokenKey(appName string) ([]byte, error) {
	if appName == "" {
		appName = "skinprice"
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve user config dir: %w", err)
	}
	keysDir := filepath.Join(configDir, appName, "keys")
	if err = os.MkdirAll(keysDir, 0o700); err != nil {
		return nil, fmt.Errorf("create key directory: %w", err)
	}
	keyPath := filepath.Join(keysDir, tokenKeyFileName)

	data, err := os.ReadFile(keyPath)
	if err == nil {
		if len(data) != 32 {
			return nil, fmt.Errorf("invalid token key length: got %d bytes, expected 32 bytes", len(data))
		}
		return data, nil
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read token key file: %w", err)
	}

	key := make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		return nil, fmt.Errorf("generate token key: %w", err)
	}
	if err = os.WriteFile(keyPath, key, 0o600); err != nil {
		return nil, fmt.Errorf("write token key file: %w", err)
	}
	return key, nil
}
