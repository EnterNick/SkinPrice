package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOrCreateTokenKey(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	keyA, err := LoadOrCreateTokenKey("skinprice-test")
	if err != nil {
		t.Fatalf("LoadOrCreateTokenKey() create error = %v", err)
	}
	if len(keyA) != 32 {
		t.Fatalf("unexpected key length: %d", len(keyA))
	}

	keyB, err := LoadOrCreateTokenKey("skinprice-test")
	if err != nil {
		t.Fatalf("LoadOrCreateTokenKey() load error = %v", err)
	}
	if string(keyA) != string(keyB) {
		t.Fatal("expected same key on repeated load")
	}

	keyPath := filepath.Join(tmp, "skinprice-test", "keys", tokenKeyFileName)
	if _, err = os.Stat(keyPath); err != nil {
		t.Fatalf("key file must exist: %v", err)
	}
}
