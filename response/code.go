package response

import "errors"

const (
	CodeSuccess    = "SUS0000"
	CodeBadRequest = "ERR0001"
)

const (
	MsgSuccess = "Thành công"
)

var (
	// ErrInvalidKey indicates that the encryption key is unreadable
	ErrInvalidKey = errors.New("invalid encryption key")

	// ErrShortCiphertext indicates that the ciphertext is too short
	ErrShortCipher = errors.New("ciphertext too short")
)
