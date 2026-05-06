package str

import (
	"fmt"
)

func ClaimStr(claims map[string]any, key string) string {
	if claims == nil {
		return ""
	}
	v := claims[key]
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	default:
		return fmt.Sprint(x)
	}
}
