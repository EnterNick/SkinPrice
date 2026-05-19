package skins

import "errors"

var (
	ErrNewSkinsRequestFailed      = errors.New("new skins request failed")
	ErrNewSkinsRequestBadStatus   = errors.New("new skins request returned bad status")
	ErrNewSkinsResponseDecodeFail = errors.New("new skins response decode failed")
	ErrNewSkinsResponseUnsuccess  = errors.New("new skins response unsuccessful")
	ErrLisSkinsTokenMissing       = errors.New("lisskins token missing")
	ErrLisSkinsTokenInvalid       = errors.New("lisskins token invalid")
)
