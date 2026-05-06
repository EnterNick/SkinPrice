package utils

import (
	"os"
	"strconv"
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
