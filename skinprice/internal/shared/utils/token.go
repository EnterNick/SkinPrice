package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateTokenByBase64(length int) (*string, error) {
	encoding := base64.StdEncoding.WithPadding(base64.NoPadding)
	out := make([]byte, length*2)
	_, err := rand.Read(out)
	if err != nil {
		return nil, err
	}
	token := encoding.EncodeToString(out)[:length]
	return &token, nil
}
