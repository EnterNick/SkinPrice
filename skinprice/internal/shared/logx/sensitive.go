package logx

import (
	"fmt"
	"strings"
)

var sensitiveKeyParts = []string{"token", "authorization", "apikey", "api_key", "secret", "password"}

func SanitizeAttrs(attrs map[string]any) map[string]any {
	if len(attrs) == 0 {
		return attrs
	}
	out := make(map[string]any, len(attrs))
	for k, v := range attrs {
		if isSensitiveKey(k) {
			out[k] = MaskString(fmt.Sprintf("%v", v))
			continue
		}
		out[k] = v
	}
	return out
}

func MaskString(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 6 {
		return "***"
	}
	return value[:3] + "***" + value[len(value)-3:]
}

func isSensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.TrimSpace(key))
	for _, part := range sensitiveKeyParts {
		if strings.Contains(normalized, part) {
			return true
		}
	}
	return false
}
