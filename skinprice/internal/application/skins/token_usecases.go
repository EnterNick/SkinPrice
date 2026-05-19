package skins

import (
	"SkinPrice/skinprice/internal/shared/errx"
	"strings"
)

const maxLisSkinsTokenLength = 1024

type TokenCipher interface {
	Encrypt(plain string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type SaveLisSkinsToken struct {
	Storage LisSkinsTokenStorage
	Cipher  TokenCipher
}

func (uc SaveLisSkinsToken) Execute(token string) error {
	normalized, err := normalizeLisSkinsToken(token)
	if err != nil {
		return err
	}
	encrypted, err := uc.Cipher.Encrypt(normalized)
	if err != nil {
		return errx.E("save_lisskins_token.encrypt", errx.CodeInternal, "failed to save lisskins token", err)
	}
	if err = uc.Storage.UpsertLisSkinsToken(encrypted); err != nil {
		return errx.E("save_lisskins_token.upsert", errx.CodeInternal, "failed to save lisskins token", err)
	}
	return nil
}

type GetLisSkinsToken struct {
	Storage LisSkinsTokenStorage
	Cipher  TokenCipher
}

func (uc GetLisSkinsToken) Execute() (string, error) {
	encrypted, err := uc.Storage.GetLisSkinsToken()
	if err != nil {
		return "", err
	}
	decrypted, err := uc.Cipher.Decrypt(encrypted)
	if err != nil {
		return "", errx.E("get_lisskins_token.decrypt", errx.CodeInternal, "failed to load lisskins token", err)
	}
	return decrypted, nil
}

type HasLisSkinsToken struct {
	Storage LisSkinsTokenStorage
}

func (uc HasLisSkinsToken) Execute() (bool, error) {
	_, err := uc.Storage.GetLisSkinsToken()
	if err == nil {
		return true, nil
	}
	if err == ErrLisSkinsTokenMissing {
		return false, nil
	}
	return false, errx.E("has_lisskins_token.get", errx.CodeInternal, "failed to check lisskins token", err)
}

type ClearLisSkinsToken struct {
	Storage LisSkinsTokenStorage
}

func (uc ClearLisSkinsToken) Execute() error {
	if err := uc.Storage.DeleteLisSkinsToken(); err != nil {
		return errx.E("clear_lisskins_token.delete", errx.CodeInternal, "failed to clear lisskins token", err)
	}
	return nil
}

func normalizeLisSkinsToken(token string) (string, error) {
	normalized := strings.TrimSpace(token)
	if normalized == "" {
		return "", errx.E("normalize_lisskins_token.empty", errx.CodeInvalidArgument, ErrLisSkinsTokenInvalid.Error(), ErrLisSkinsTokenInvalid)
	}
	if len(normalized) > maxLisSkinsTokenLength {
		return "", errx.E("normalize_lisskins_token.length", errx.CodeInvalidArgument, ErrLisSkinsTokenInvalid.Error(), ErrLisSkinsTokenInvalid)
	}
	return normalized, nil
}
