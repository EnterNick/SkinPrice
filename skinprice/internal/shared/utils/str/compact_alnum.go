package str

import (
	"strings"
	"unicode"
)

func CompactAlnum(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if unicode.IsLetter(r) {
			b.WriteRune(r)
			continue
		}
		if unicode.IsNumber(r) {
			if mapped, ok := leetMap(r); ok {
				b.WriteRune(mapped)
				continue
			}
			b.WriteRune(r)
			continue
		}
		if mapped, ok := leetMap(r); ok {
			b.WriteRune(mapped)
			continue
		}
	}

	return b.String()
}

func leetMap(r rune) (rune, bool) {
	switch r {
	case '0':
		return 'o', true
	case '1':
		return 'i', true
	case '3':
		return 'e', true
	case '4':
		return 'a', true
	case '5':
		return 's', true
	case '7':
		return 't', true
	case '@':
		return 'a', true
	case '$':
		return 's', true
	default:
		return 0, false
	}
}
