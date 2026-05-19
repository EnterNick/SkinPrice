package skins

type TokenCipher interface {
	Encrypt(plain string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type SaveLisSkinsToken struct {
	Storage LisSkinsTokenStorage
	Cipher  TokenCipher
}

func (uc SaveLisSkinsToken) Execute(token string) error {
	encrypted, err := uc.Cipher.Encrypt(token)
	if err != nil {
		return err
	}
	return uc.Storage.UpsertLisSkinsToken(encrypted)
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
	return uc.Cipher.Decrypt(encrypted)
}
