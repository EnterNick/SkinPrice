package settings

import "testing"

func TestNormalizeFontFamily(t *testing.T) {
	t.Run("keeps supported values", func(t *testing.T) {
		for _, value := range []string{"inter", "system", "nunito", "roboto", "ibm-plex-sans", "manrope", "monocraft"} {
			if actual := normalizeFontFamily(value); actual != value {
				t.Fatalf("expected %q, got %q", value, actual)
			}
		}
	})

	t.Run("falls back to default", func(t *testing.T) {
		if actual := normalizeFontFamily(""); actual != DefaultFontFamily {
			t.Fatalf("expected default %q, got %q", DefaultFontFamily, actual)
		}
		if actual := normalizeFontFamily("serif"); actual != DefaultFontFamily {
			t.Fatalf("expected default %q, got %q", DefaultFontFamily, actual)
		}
	})
}

func TestNormalizeFontSizePx(t *testing.T) {
	t.Run("falls back to default", func(t *testing.T) {
		if actual := normalizeFontSizePx(0); actual != DefaultFontSizePx {
			t.Fatalf("expected default %d, got %d", DefaultFontSizePx, actual)
		}
		if actual := normalizeFontSizePx(100); actual != DefaultFontSizePx {
			t.Fatalf("expected default %d, got %d", DefaultFontSizePx, actual)
		}
	})

	t.Run("keeps allowed range", func(t *testing.T) {
		for _, value := range []int{10, 14, 18, 28} {
			if actual := normalizeFontSizePx(value); actual != value {
				t.Fatalf("expected %d, got %d", value, actual)
			}
		}
	})
}

func TestNormalizeAutoRefreshEnabled(t *testing.T) {
	if !normalizeAutoRefreshEnabled(true) {
		t.Fatal("expected true to remain true")
	}
	if normalizeAutoRefreshEnabled(false) {
		t.Fatal("expected false to remain false")
	}
}
