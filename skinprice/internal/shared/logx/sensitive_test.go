package logx

import "testing"

func TestSanitizeAttrsMasksSensitiveKeys(t *testing.T) {
	attrs := map[string]any{
		"token":         "abcdefghijklmnopqrstuvwxyz",
		"Authorization": "Bearer super-secret-token",
		"market_hash":   "AK-47",
	}
	sanitized := SanitizeAttrs(attrs)
	if sanitized["token"] == attrs["token"] {
		t.Fatalf("expected token to be masked")
	}
	if sanitized["Authorization"] == attrs["Authorization"] {
		t.Fatalf("expected Authorization to be masked")
	}
	if sanitized["market_hash"] != attrs["market_hash"] {
		t.Fatalf("expected non-sensitive attr unchanged")
	}
}
