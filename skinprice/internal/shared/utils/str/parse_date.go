package str

import (
	"strings"
	"time"
)

func ParseBirthDate(v string) time.Time {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}
	}
	if t, err := time.Parse("2006-01-02", v); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return t
	}
	return time.Time{}
}
