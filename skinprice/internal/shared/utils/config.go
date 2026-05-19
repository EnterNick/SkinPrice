package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetStrWDefault(name string, def string) string {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	return value
}

func GetIntWDefault(name string, def int) int {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return valueInt
}

func GetBoolWDefault(name string, def bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(name)))
	if value == "" {
		return def
	}

	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return def
	}
}
