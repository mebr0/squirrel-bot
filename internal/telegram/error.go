package telegram

import "errors"

var (
	ErrUnknownCallback = errors.New("unsupported callback")
)

